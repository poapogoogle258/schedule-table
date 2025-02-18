package worker

import (
	"context"
	"fmt"
	"schedule_table/internal/constant"
	"schedule_table/internal/lib"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg/logger"
	"schedule_table/internal/repository"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WorkerInit struct {
	CalRepo      repository.CalendarRepository
	TaskRepo     repository.ITaskRepository
	ScheduleRepo repository.ScheduleRepository
	MemberRepo   repository.MembersRepository
}

type TaskWorker interface {
	Start(ctx context.Context) error
}

type taskWorker struct {
	Id           int
	JobQueue     *JobQueue
	CalRepo      repository.CalendarRepository
	TaskRepo     repository.ITaskRepository
	ScheduleRepo repository.ScheduleRepository
	MemberRepo   repository.MembersRepository

	calendar      *dao.Calendars
	employee      map[uuid.UUID]lib.IWorker
	employeeQueue *lib.WorkerQueueMap
}

func selectWorkersResponsible(workers map[uuid.UUID]lib.IWorker, responsibly []*dao.Responsible) []lib.IWorker {
	selectWorkers := make([]lib.IWorker, 0)

	for _, responsible := range responsibly {
		if _, ok := workers[responsible.MemberId]; ok {
			selectWorkers = append(selectWorkers, workers[responsible.MemberId])
		}
	}

	return selectWorkers
}

func softByDateTimeAndPriority(a, b *dao.Tasks) int {
	if c := a.Start.Compare(b.Start); c == 0 {
		if a.Priority > b.Priority {
			return 1
		} else {
			return -1
		}
	} else {
		return c
	}
}

func initializeWorkers(members []*dao.Members) map[uuid.UUID]lib.IWorker {
	workers := make(map[uuid.UUID]lib.IWorker)

	for i := range members {
		workers[members[i].Id] = lib.NewWorkerMember(members[i])
	}

	return workers
}

func getRangeTimeToGenerateTask() (start time.Time, end time.Time) {
	today := time.Now().Local()
	start = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()) // today 00:00:00
	end = start.AddDate(0, 1, -start.Day())                                                   // end of month

	return
}

func (worker *taskWorker) generateTasks(start time.Time, end time.Time) ([]*dao.Tasks, []*dao.Tasks, error) {

	tasks := make([]*dao.Tasks, 0)
	tasksDifference := make([]*dao.Tasks, 0)
	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup
		for i := range worker.calendar.Schedules {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				schedule := worker.calendar.Schedules[i]

				if schedule.MasterScheduleId == nil {
					worker.employeeQueue.Set(schedule.Id, lib.NewWorkerQueue(selectWorkersResponsible(worker.employee, schedule.Responsibles)))
				}

				today := time.Now().Local()
				startOfCreateTasks := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
				tasksGenerated := lib.CreateRecurrenceTasks(schedule, startOfCreateTasks, end)
				tasksCalendar, errQueryTaskCalendar := worker.TaskRepo.FindWithAssociation("schedule_id = ? AND start BETWEEN ? AND ?", schedule.Id, start, end)
				if errQueryTaskCalendar != nil {
					errCh <- errQueryTaskCalendar

					return
				}

				// marge taskGenerated and taskCalendar
				for _, taskCalendar := range tasksCalendar {

					// Check Revered Task and add task to worker
					if taskCalendar.Status == constant.TaskStatus_Reserved && taskCalendar.MemberId != nil {
						if _, ok := worker.employee[*taskCalendar.MemberId]; ok {
							worker.employee[*taskCalendar.MemberId].AddReservedTask(taskCalendar)
						}
					}

					// check RecurrenceId in taskCalendars
					IndexGenerated := slices.IndexFunc(tasksGenerated, func(taskGenerated *dao.Tasks) bool {
						return taskGenerated.RecurrenceId == taskCalendar.RecurrenceId
					})

					if IndexGenerated != -1 {
						// instead of tasksGenerated with taskCalendar
						tasksGenerated[IndexGenerated] = taskCalendar
					} else {
						// sum task not in tasksGenerated
						tasksDifference = append(tasksDifference, taskCalendar)
					}

				}

				tasks = append(tasks, tasksGenerated...)
			}(i)
		}
		wg.Wait()
		close(errCh)
	}()

	if err, haveErr := <-errCh; haveErr {
		return nil, nil, err
	} else {
		// sort tasks
		slices.SortFunc(tasks, softByDateTimeAndPriority)
		return tasks, tasksDifference, nil
	}

}

func (worker *taskWorker) orderWorkerQueue(start time.Time) error {

	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup
		for masterId := range worker.employeeQueue.Items {
			wg.Add(1)
			go func(masterId uuid.UUID) {
				defer wg.Done()
				schedulesId := make([]uuid.UUID, 0)
				for _, schedule := range worker.calendar.Schedules {
					if schedule.Id == masterId || (schedule.MasterScheduleId != nil && *schedule.MasterScheduleId == masterId) {
						schedulesId = append(schedulesId, schedule.Id)
					}
				}

				tasksCalendar, errQueryTaskCalendar := worker.TaskRepo.FindOrderLimit("start DESC", worker.employeeQueue.Get(masterId).Size(), "schedule_id IN ? AND start < ? AND status != ?", schedulesId, start, constant.TaskStatus_Canceled)
				if errQueryTaskCalendar != nil {
					errCh <- errQueryTaskCalendar
					return
				}

				worker.employeeQueue.Get(masterId).OrderQueue(tasksCalendar)
			}(masterId)
		}
		wg.Wait()
		close(errCh)
	}()

	if err, haveErr := <-errCh; haveErr {
		return err
	}

	return nil
}

func (worker *taskWorker) Start(ctx context.Context) error {
	logger.Debug("start:UpdateTasksOfCalendarIdWorker", zap.Int("id", worker.Id))
	defer logger.Debug("end:UpdateTasksOfCalendarIdWorker", zap.Int("id", worker.Id))

	for {
		select {
		case <-ctx.Done():
			return nil
		case message := <-worker.JobQueue.Queue:
			prefixLog := fmt.Sprintf("taskWorker[%d] jobId: %s ", worker.Id, message.Id)

			logger.Debug(prefixLog+"processing", zap.String("calendarId", message.CalendarId), zap.Time("updateAt", message.UpdatedAt))

			calendarId := message.CalendarId

			start, end := getRangeTimeToGenerateTask()

			// set calendar
			if calendar, err := worker.CalRepo.FindOneWithAssociation(calendarId, start, end); err != nil {
				logger.Error(prefixLog+"findOneWithAssociationError", zap.Error(err))
				continue
			} else {
				worker.calendar = calendar
			}

			// set Member To Employee
			worker.employee = initializeWorkers(worker.calendar.Members)

			// set Employee Queue
			worker.employeeQueue = lib.NewWorkerQueueMap()
			if err := worker.orderWorkerQueue(start); err != nil {
				logger.Error(prefixLog+"orderWorkerQueueError", zap.Error(err))
				continue
			}

			// get tasks
			tasks, tasksDifference, errGenerateTasks := worker.generateTasks(start, end)
			if errGenerateTasks != nil {
				logger.Error(prefixLog+"generateTasksError", zap.Error(errGenerateTasks))
				continue
			}

			// match tasks
			for i := 0; i < len(tasks); i++ {
				var selectWorkerQueueId uuid.UUID
				if tasks[i].Description.MasterScheduleId == nil {
					selectWorkerQueueId = tasks[i].Description.Id
				} else {
					selectWorkerQueueId = *tasks[i].Description.MasterScheduleId
				}

				worker.employeeQueue.Get(selectWorkerQueueId).Match(tasks[i])
			}

			// save to database
			lastTimeTaskUpdated, errGetLastTimeTaskUpdatedOfCalendarId := worker.CalRepo.GetLastTimeGeneratedTask(calendarId)
			if errGetLastTimeTaskUpdatedOfCalendarId != nil {
				logger.Error(prefixLog+"getLastTimeGeneratedTaskError", zap.Error(errGetLastTimeTaskUpdatedOfCalendarId))
				continue
			}

			if lastTimeTaskUpdated == nil || lastTimeTaskUpdated.After(message.UpdatedAt) {
				// CleanUp tasksDifference
				worker.TaskRepo.Delete(tasksDifference)

				// Upsert tasks
				worker.TaskRepo.Upsert(tasks)

				// update lastTime Task updated
				worker.CalRepo.UpdateTaskGenerated(calendarId, message.UpdatedAt)

				logger.Debug(prefixLog + "success")
			}
		}
	}
}

func NewTaskWorker(id int, jobQueue *JobQueue, workerInit *WorkerInit) TaskWorker {

	return &taskWorker{
		Id:           id,
		JobQueue:     jobQueue,
		CalRepo:      workerInit.CalRepo,
		TaskRepo:     workerInit.TaskRepo,
		ScheduleRepo: workerInit.ScheduleRepo,
		MemberRepo:   workerInit.MemberRepo,
	}
}
