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
	"github.com/jinzhu/copier"
)

type TasksHandler interface {
	GetTasks(c *gin.Context) (*[]dto.ResponseTask, error)
}

type tasksHandler struct {
	CalRepo     repository.CalendarRepository
	ScheService service.IScheduleService
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
	if !taskHandler.CalRepo.IsExits(calendarId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("not found calendar id"))
	}

	start := util.Must(time.Parse(time.RFC3339, query.Start))
	end := util.Must(time.Parse(time.RFC3339, query.End))

	calendar, errFindCalendar := taskHandler.CalRepo.FindOneWithAssociation(calendarId, start, end)
	if errFindCalendar != nil {
		return nil, errFindCalendar
	}

	calendarTasks := make([]dao.Tasks, 0)

	for i := 0; i < len(*calendar.Schedules); i++ {
		schedule := (*calendar.Schedules)[i]
		tasks := taskHandler.ScheService.NewSchedule(&schedule).GenerateTasks(start, end)
		calendarTasks = append(calendarTasks, (*tasks)...)
	}

	slices.SortFunc(calendarTasks, softByDateTimeAndPriority)

	response := &[]dto.ResponseTask{}
	if err := copier.Copy(&response, &calendarTasks); err != nil {
		return nil, err
	}

	return response, nil
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

func NewTasksHandler(calRepo repository.CalendarRepository, scheService service.IScheduleService) TasksHandler {
	return &tasksHandler{
		CalRepo:     calRepo,
		ScheService: scheService,
	}
}
