package handler

import (
	"errors"
	"fmt"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/lib"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TasksHandler interface {
	GetTasks(c *gin.Context) (*[]dto.ResponseTask, error)
	EditTask(c *gin.Context) (*dto.ResponseTask, error)
}

type tasksHandler struct {
	CalRepo      repository.CalendarRepository
	TaskRepo     repository.ITaskRepository
	ScheduleRepo repository.ScheduleRepository
	MemberRepo   repository.MembersRepository
}

type queryStringGetTasks struct {
	Start string `form:"start" binding:"required"`
	End   string `form:"end" binding:"required"`
}

func (taskHandler *tasksHandler) GetTasks(c *gin.Context) (*[]dto.ResponseTask, error) {

	var query queryStringGetTasks
	if err := c.BindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("query string not validate"))
	}

	calendarId := c.Param("calendarId")
	if err := taskHandler.CalRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	start := util.Must(time.Parse(time.RFC3339, query.Start))
	end := util.Must(time.Parse(time.RFC3339, query.End))

	if !taskHandler.CalRepo.CheckRecurrenceChanged(calendarId) {
		if tasks, err := taskHandler.TaskRepo.FindWithAssociation("calendar_id = ? AND start BETWEEN ? AND ?", calendarId, start, end); err != nil {
			return nil, err
		} else {
			return util.Convert[[]dto.ResponseTask](&tasks), nil
		}
	}

	// ------------------ generate new task --------------------------------

	calendar, errFindCalendar := taskHandler.CalRepo.FindOneWithAssociation(calendarId, start, end)
	if errFindCalendar != nil {
		return nil, errFindCalendar
	}

	workers := make(map[uuid.UUID]lib.IWorker)
	workerQueue := make(map[uuid.UUID]lib.IWorkerQueue)
	tasks := make([]*dao.Tasks, 0)
	tasksDifference := make([]*dao.Tasks, 0)

	var workerQueueMutex sync.Mutex

	// create worker
	for _, member := range *calendar.Members {
		workers[member.Id] = lib.NewWorkerMember(&member)
	}

	// create tasks
	var wg sync.WaitGroup
	for i := 0; i < len(*calendar.Schedules); i++ {
		wg.Add(1)
		go func() {
			schedule := (*calendar.Schedules)[i]

			if schedule.MasterScheduleId == nil {
				// workerQueue[schedule.Id] = lib.NewWorkerQueue(selectWorkersResponsible(workers, *schedule.Responsibles))
				workerQueueMutex.Lock()
				workerQueue[schedule.Id] = lib.NewWorkerQueue(selectWorkersResponsible(workers, *schedule.Responsibles))
				workerQueueMutex.Unlock()
			}
			tasksGenerated := lib.CreateRecurrenceTasks(&schedule, start, end)
			tasksCalendar, errQueryTaskCalendar := taskHandler.TaskRepo.FindWithAssociation("schedule_id = ? AND start BETWEEN ? AND ?", schedule.Id, start, end)
			if errQueryTaskCalendar != nil {
				panic(errQueryTaskCalendar)
			}

			// marge taskGenerated and taskCalendar
		TaskCalendarLoop:
			for _, taskCalendar := range *tasksCalendar {

				// add Revered Task to worker
				if taskCalendar.Status == constant.TaskStatus_Reserved && taskCalendar.MemberId != nil {
					if _, ok := workers[*taskCalendar.MemberId]; ok {
						workers[*taskCalendar.MemberId].AddReservedTask(&taskCalendar)
					}
				}

				// check marge
				for j, taskGenerated := range tasksGenerated {
					if taskGenerated.RecurrenceId == taskCalendar.RecurrenceId {
						tasksGenerated[j] = &taskCalendar
						continue TaskCalendarLoop
					}
				}

				tasksDifference = append(tasksDifference, &taskCalendar)
			}

			tasks = append(tasks, tasksGenerated...)
			wg.Done()
		}()
	}
	wg.Wait()

	// handler tasksDifference
	for _, task := range tasksDifference {
		go taskHandler.TaskRepo.DeleteOne(task.Id)
	}

	// order queue worker
	for masterId := range workerQueue {
		wg.Add(1)
		go func() {
			schedulesId := make([]uuid.UUID, 0)
			for _, schedule := range *calendar.Schedules {
				if schedule.Id == masterId || (schedule.MasterScheduleId != nil && *schedule.MasterScheduleId == masterId) {
					schedulesId = append(schedulesId, schedule.Id)
				}
			}

			tasksCalendar, errQueryTaskCalendar := taskHandler.TaskRepo.FindOrderLimit("start DESC", workerQueue[masterId].Size(), "schedule_id IN ? AND start < ? AND status != ?", schedulesId, start, constant.TaskStatus_Canceled)
			if errQueryTaskCalendar != nil {
				panic(errQueryTaskCalendar)
			}

			workerQueue[masterId].OrderQueue(tasksCalendar)
			wg.Done()
		}()
	}
	wg.Wait()

	// sort tasks
	slices.SortFunc(tasks, softByDateTimeAndPriority)

	// match tasks
	for i := 0; i < len(tasks); i++ {
		var selectWorkerQueueId uuid.UUID
		if tasks[i].Description.MasterScheduleId == nil {
			selectWorkerQueueId = tasks[i].Description.Id
		} else {
			selectWorkerQueueId = *tasks[i].Description.MasterScheduleId
		}

		if selectWorkerQueueId == uuid.Nil {
			fmt.Println("tasks", tasks[i])
			return nil, errors.New("selectWorkerQueueId is nil")
		}

		if _, checkQueue := workerQueue[selectWorkerQueueId]; !checkQueue {
			return nil, errors.New("not have workerQueue")
		}

		workerQueue[selectWorkerQueueId].Match(tasks[i])
	}

	taskHandler.TaskRepo.CreateTasks(tasks)
	taskHandler.CalRepo.UpdateLastTimeGenerated(calendarId)

	return util.Convert[[]dto.ResponseTask](&tasks), nil
}

func selectWorkersResponsible(workers map[uuid.UUID]lib.IWorker, responsibly []dao.Responsible) []lib.IWorker {
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

type ReserveMemberBody struct {
	MemberId *uuid.UUID `json:"member_id" binding:"-"`
	Start    time.Time  `json:"start" binding:"required"`
	End      time.Time  `json:"end" binding:"required"`
	Status   int8       `json:"status" binding:"required"`
}

var (
	ErrStartDateTimeIsAfterEndDateTime = errors.New("start datetime is after end datetime")
)

func (handler *tasksHandler) EditTask(c *gin.Context) (*dto.ResponseTask, error) {

	calendarId := c.Param("calendarId")
	if !handler.CalRepo.IsExist(calendarId) {
		return nil, pkg.NewErrorWithStatusCode(404, repository.ErrCalendarNotFount)
	}

	taskId := c.Param("taskId")
	if !handler.TaskRepo.IsExist(taskId) {
		return nil, pkg.NewErrorWithStatusCode(404, repository.ErrTaskNotExists)
	}

	body := &ReserveMemberBody{}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	if body.Start.After(body.End) {
		return nil, pkg.NewErrorWithStatusCode(400, ErrStartDateTimeIsAfterEndDateTime)
	}

	if _, err := constant.NewTaskStatus(body.Status); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	data := map[string]interface{}{
		"member_id": body.MemberId,
		"start":     body.Start,
		"end":       body.End,
		"status":    constant.TaskStatus(body.Status),
	}

	task, errUpdateTask := handler.TaskRepo.UpdatesAndFind(taskId, data)
	if errUpdateTask != nil {
		return nil, pkg.NewErrorWithStatusCode(500, errUpdateTask)
	}

	return util.Convert[dto.ResponseTask](&task), nil
}

func NewTasksHandler(calRepo repository.CalendarRepository, taskRepo repository.ITaskRepository, scheduleRepo repository.ScheduleRepository, memberRepo repository.MembersRepository) TasksHandler {

	return &tasksHandler{
		CalRepo:      calRepo,
		TaskRepo:     taskRepo,
		ScheduleRepo: scheduleRepo,
		MemberRepo:   memberRepo,
	}
}
