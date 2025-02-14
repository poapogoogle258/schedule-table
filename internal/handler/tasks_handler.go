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
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var (
	ErrQueryStartMustFormat = errors.New("format start query string must be RFC3339 2006-01-02T15:04:05%2b+07:00 (encode '+' to ASCII %2b)")
	ErrQueryEndMustFormat   = errors.New("format end query string must be RFC3339 2006-01-02T15:04:05%2b+07:00 (encode '+' to ASCII %2b)")
)

type TasksHandler interface {
	GetTasks(c *gin.Context) ([]*dto.ResponseTask, error)
	EditTask(c *gin.Context) (*dto.ResponseTask, error)
}

type tasksHandler struct {
	CalRepo      repository.CalendarRepository
	TaskRepo     repository.ITaskRepository
	ScheduleRepo repository.ScheduleRepository
	MemberRepo   repository.MembersRepository
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

type queryStringGetTasks struct {
	Start string `form:"start" binding:"required"`
	End   string `form:"end" binding:"required"`
}

func (taskHandler *tasksHandler) GetTasks(c *gin.Context) ([]*dto.ResponseTask, error) {
	calendarId := c.Param("calendarId")

	var query queryStringGetTasks
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	start, errParseStart := time.Parse(time.RFC3339, query.Start)
	if errParseStart != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errParseStart)
	}
	end, errParseEnd := time.Parse(time.RFC3339, query.End)
	if errParseEnd != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrQueryEndMustFormat)
	}

	if !taskHandler.CalRepo.CheckRecurrenceChanged(calendarId) {
		tasks, errFindWithAssociation := taskHandler.TaskRepo.FindWithAssociation("calendar_id = ? AND start BETWEEN ? AND ?", calendarId, query.Start, query.End)
		if errFindWithAssociation != nil {
			return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errFindWithAssociation)
		}

		resp := []*dto.ResponseTask{}
		copier.Copy(&resp, tasks)

		return resp, nil
	}

	// ------------------ generate new task --------------------------------

	calendar, errFindOneWithAssociation := taskHandler.CalRepo.FindOneWithAssociation(calendarId, start, end)
	if errFindOneWithAssociation != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errFindOneWithAssociation)
	}

	workers := make(map[uuid.UUID]lib.IWorker)
	workerQueue := make(map[uuid.UUID]lib.IWorkerQueue)
	tasks := make([]*dao.Tasks, 0)
	tasksDifference := make([]*dao.Tasks, 0)

	// var workerQueueMutex sync.Mutex
	// workerQueueMutex.Unlock()

	// create worker
	for _, member := range *calendar.Members {
		workers[member.Id] = lib.NewWorkerMember(&member)
	}

	// create tasks
	var wg sync.WaitGroup
	for i := range *calendar.Schedules {
		wg.Add(1)
		go func() {
			defer wg.Done()

			schedule := (*calendar.Schedules)[i]

			if schedule.MasterScheduleId == nil {
				workerQueue[schedule.Id] = lib.NewWorkerQueue(selectWorkersResponsible(workers, *schedule.Responsibles))
			}
			tasksGenerated := lib.CreateRecurrenceTasks(&schedule, start, end)
			tasksCalendar, errQueryTaskCalendar := taskHandler.TaskRepo.FindWithAssociation("schedule_id = ? AND start BETWEEN ? AND ?", schedule.Id, start, end)
			if errQueryTaskCalendar != nil {
				panic(errQueryTaskCalendar) // fix this
			}

			// marge taskGenerated and taskCalendar
		TaskCalendarLoop:
			for _, taskCalendar := range tasksCalendar {

				// add Revered Task to worker
				if taskCalendar.Status == constant.TaskStatus_Reserved && taskCalendar.MemberId != nil {
					if _, ok := workers[*taskCalendar.MemberId]; ok {
						workers[*taskCalendar.MemberId].AddReservedTask(taskCalendar)
					}
				}

				// check marge
				for j, taskGenerated := range tasksGenerated {
					if taskGenerated.RecurrenceId == taskCalendar.RecurrenceId {
						tasksGenerated[j] = taskCalendar
						continue TaskCalendarLoop
					}
				}

				tasksDifference = append(tasksDifference, taskCalendar)
			}

			tasks = append(tasks, tasksGenerated...)
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

	// taskHandler.TaskRepo.CreateTasks(tasks) // TO DO: fix this to CreateOrUpdate row
	taskHandler.CalRepo.UpdateLastTimeGenerated(calendarId)

	resp := []*dto.ResponseTask{}
	copier.Copy(&resp, tasks)

	return resp, nil
}

var (
	ErrStartDateTimeIsAfterEndDateTime = errors.New("start datetime is after end datetime")
)

type ReserveMemberBody struct {
	MemberId *uuid.UUID `json:"member_id" binding:"-"`
	Start    time.Time  `json:"start" binding:"required"`
	End      time.Time  `json:"end" binding:"required"`
	Status   int8       `json:"status" binding:"required"`
}

func (body *ReserveMemberBody) Validate() error {
	if body.Start.After(body.End) {
		return ErrStartDateTimeIsAfterEndDateTime
	}

	if _, err := constant.NewTaskStatus(body.Status); err != nil {
		return err
	}

	return nil
}

func (handler *tasksHandler) EditTask(c *gin.Context) (*dto.ResponseTask, error) {

	taskId := c.Param("taskId")

	body := &ReserveMemberBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	if err := body.Validate(); err != nil {
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
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errUpdateTask)
	}

	resp := dto.ResponseTask{}
	copier.Copy(&resp, task)

	return &resp, nil

}

func NewTasksHandler(calRepo repository.CalendarRepository, taskRepo repository.ITaskRepository, scheduleRepo repository.ScheduleRepository, memberRepo repository.MembersRepository) TasksHandler {

	return &tasksHandler{
		CalRepo:      calRepo,
		TaskRepo:     taskRepo,
		ScheduleRepo: scheduleRepo,
		MemberRepo:   memberRepo,
	}
}
