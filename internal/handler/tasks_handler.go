package handler

import (
	"errors"
	"net/http"
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
	ReserveMember(c *gin.Context) (*dao.Tasks, error)
}

type tasksHandler struct {
	CalRepo  repository.CalendarRepository
	TaskRepo repository.ITaskRepository
}

type queryStringGetTasks struct {
	Start  string `form:"start" binding:"required"`
	End    string `form:"end" binding:"required"`
	Action string `form:"action" binding:"required"`
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

	calendar, errFindCalendar := taskHandler.CalRepo.FindOneWithAssociation(calendarId, start, end)
	if errFindCalendar != nil {
		return nil, errFindCalendar
	}

	workers := make([]lib.IWorker, 0)
	workerQueue := make(map[uuid.UUID]lib.IWorkerQueue)
	tasks := make([]*dao.Tasks, 0)
	tasksDifference := make([]*dao.Tasks, 0)

	// create worker
	for _, member := range *calendar.Members {
		workers = append(workers, lib.NewWorkerMember(&member))
	}

	// create tasks
	var wg sync.WaitGroup
	for i := 0; i < len(*calendar.Schedules); i++ {
		wg.Add(1)
		go func() {
			schedule := (*calendar.Schedules)[i]

			if schedule.MasterScheduleId == nil {
				workerQueue[schedule.Id] = lib.NewWorkerQueue(selectWorkersResponsible(&workers, schedule.Responsibles))
			}
			tasksGenerated := lib.CreateRecurrenceTasks(&schedule, start, end)
			tasksCalendar, errQueryTaskCalendar := taskHandler.TaskRepo.Find("schedule_id = ? AND start BETWEEN ? AND ?", schedule.Id, start, end)
			if errQueryTaskCalendar != nil {
				panic(errQueryTaskCalendar)
			}

			// marge taskGenerated and taskCalendar to { Marge, Difference }
		TaskCalendarLoop:
			for _, taskCalendar := range *tasksCalendar {
				for j, taskGenerated := range tasksGenerated {
					if taskGenerated.Start.Equal(taskCalendar.Start) {
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

			tasksCalendar, errQueryTaskCalendar := taskHandler.TaskRepo.FindOrderLimit("start DESC", workerQueue[masterId].Size(), "schedule_id IN ? AND start BETWEEN ? AND ?", schedulesId, start, end)
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

		if _, checkQueue := workerQueue[selectWorkerQueueId]; !checkQueue {
			return nil, errors.New("not have workerQueue")
		}

		workerQueue[selectWorkerQueueId].Match(tasks[i])
	}

	return util.Convert[[]dto.ResponseTask](&tasks), nil
}

func selectWorkersResponsible(workers *[]lib.IWorker, responsibly *[]dao.Responsible) []lib.IWorker {
	selectWorkers := make([]lib.IWorker, 0)

	for _, responsible := range *responsibly {
		for i := 0; i < len(*workers); i++ {
			if responsible.MemberId == (*workers)[i].GetId() {
				selectWorkers = append(selectWorkers, (*workers)[i])
				break
			}
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

type ReserveMemberQueryString struct {
	MemberId string `form:"member_id"`
	Reserved string `form:"reserved"`
}

func (handler *tasksHandler) ReserveMember(c *gin.Context) (*dao.Tasks, error) {
	taskId := c.Param("taskId")
	memberId := c.Param("memberId")

	var query ReserveMemberQueryString
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("bad request Must have 'member_id' and 'reserved' in query string"))
	}

	insert := map[string]interface{}{}
	if query.Reserved == "true" {
		insert["member_id"] = memberId
		insert["reserved"] = true

	} else if query.Reserved == "false" {
		insert["reserved"] = false
	} else {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("bad request in query field 'reserved' value must be 'true' or 'false'"))
	}

	task, errUpdate := handler.TaskRepo.UpdatesAndFind(taskId, insert)

	if errUpdate != nil {
		return nil, errUpdate
	}

	return task, nil
}

func NewTasksHandler(calRepo repository.CalendarRepository, taskRepo repository.ITaskRepository) TasksHandler {
	return &tasksHandler{
		CalRepo:  calRepo,
		TaskRepo: taskRepo,
	}
}
