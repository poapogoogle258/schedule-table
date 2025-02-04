package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"
	"schedule_table/util"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type TasksHandler interface {
	GetTasks(c *gin.Context) (*[]dto.ResponseTask, error)
	ReserveMember(c *gin.Context) (*dao.Tasks, error)
}

type tasksHandler struct {
	CalRepo        repository.CalendarRepository
	ScheService    service.IScheduleService
	ManagerService service.IManagerService
	TaskRepo       repository.ITaskRepository
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

	managers := make(map[uuid.UUID]*service.Manager)
	calendarTasks := make([]dao.Tasks, 0)

	for loopMasterId := 0; loopMasterId < 2; loopMasterId++ {
		for i := 0; i < len(*calendar.Schedules); i++ {
			schedule := (*calendar.Schedules)[i]

			if loopMasterId == 0 && schedule.MasterScheduleId == nil {
				manager := taskHandler.ManagerService.NewManagerSchedule(&schedule)
				manager.Tasks = taskHandler.ScheService.NewSchedule(&schedule).GenerateTasks(start, end)

				// DO TO : Refix
				for c := 0; c < len(*manager.Tasks); c++ {
					(*manager.Tasks)[c].Description = &schedule
				}
				// ----------------------------------

				managers[manager.Id] = manager
				calendarTasks = append(calendarTasks, (*manager.Tasks)...)
			} else if loopMasterId == 1 && schedule.MasterScheduleId != nil {
				masterId := *schedule.MasterScheduleId
				if _, hasManager := managers[masterId]; !hasManager {
					panic(errors.New("DO TO: load master queue to childe"))
				}

				manager := taskHandler.ManagerService.NewManagerScheduleWithQueue(&schedule, managers[masterId].Queue)
				manager.Tasks = taskHandler.ScheService.NewSchedule(&schedule).GenerateTasks(start, end)

				managers[manager.Id] = manager
				calendarTasks = append(calendarTasks, (*manager.Tasks)...)
			}
		}
	}

	slices.SortFunc(calendarTasks, softByDateTimeAndPriority)

	taskReserved, _ := taskHandler.TaskRepo.Find("(start BETWEEN @start AND @end) AND ('end' BETWEEN @start AND @end) AND reserved = true", map[string]interface{}{
		"start": start.UTC().Format(time.RFC3339),
		"end":   end.UTC().Format(time.RFC3339),
	})
	if taskReserved != nil {
		for i := 0; i < len(*taskReserved); i++ {
			for j := 0; j < len(calendarTasks); j++ {
				if checkTaskIsBooking(&(*taskReserved)[i], &calendarTasks[j]) {
					calendarTasks[j] = (*taskReserved)[i]
				}
			}
		}
	}

	for i := 0; i < len(calendarTasks); i++ {
		task := &calendarTasks[i]
		for n := 0; ; n++ {
			if err := managers[task.ScheduleId].Queue.Next(n).AddTask(task, managers[task.ScheduleId].RestTime); err == nil {
				managers[task.ScheduleId].Queue.Select(n)
				managers[task.ScheduleId].Count.Add(task.MemberId)
				break
			} else if errors.Is(err, service.ErrorSkipAllQueue) {
				// skip is task
				break
			} else {
				managers[task.ScheduleId].Queue.Skip()
				// TO DO: Handler Force
			}
		}
	}

	response := &[]dto.ResponseTask{}
	if err := copier.Copy(&response, &calendarTasks); err != nil {
		return nil, err
	}

	return response, nil
}

func checkTaskIsBooking(task, generateTask *dao.Tasks) bool {
	if generateTask.Reserved {
		return false
	}

	if task.ScheduleId == generateTask.ScheduleId && task.Start.Equal(generateTask.Start) && task.End.Equal(generateTask.End) {
		return true
	} else {
		return false
	}
}

func softByDateTimeAndPriority(a, b dao.Tasks) int {
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

func NewTasksHandler(calRepo repository.CalendarRepository, scheService service.IScheduleService, managerService service.IManagerService, taskRepo repository.ITaskRepository) TasksHandler {
	return &tasksHandler{
		CalRepo:        calRepo,
		ScheService:    scheService,
		ManagerService: managerService,
		TaskRepo:       taskRepo,
	}
}
