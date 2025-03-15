package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	GetTasksOfCalendarId(c *gin.Context) ([]*dto.TaskInfo, error)
	UpdateTasksOfCalendarId(c *gin.Context) (*dto.TaskInfo, error)
}

var ErrTaskNotFount = errors.New("task not found")

type taskHandler struct {
	taskService service.TaskService
	calService  service.CalendarService
}

type GetTasksOfCalendarsQuery struct {
	Start time.Time `form:"start" binding:"required" time_format:"2006-01-02"`
	End   time.Time `form:"end" binding:"required" time_format:"2006-01-02"`
}

func (t *taskHandler) RefreshTasksOfCalendarId(calendarId string, start time.Time, end time.Time) error {

	calendar, err := t.calService.FindCalendarIdWithFullAggregate(calendarId, start, end)
	if err != nil {
		return err
	}

	if err := t.taskService.RefreshTasksOfCalendar(calendar, start, end); err != nil {
		return err
	}

	return nil
}

func (t *taskHandler) GetTasksOfCalendarId(c *gin.Context) ([]*dto.TaskInfo, error) {

	calendarId := c.Param("calendarId")
	var query GetTasksOfCalendarsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := t.RefreshTasksOfCalendarId(calendarId, query.Start, query.End); err != nil {
		return nil, err
	}

	return t.taskService.GetTasksOfCalendars(calendarId, query.Start, query.End)
}

func (t *taskHandler) UpdateTasksOfCalendarId(c *gin.Context) (*dto.TaskInfo, error) {
	calendarId := c.Param("calendarId")
	taskId := c.Param("taskId")

	if !t.taskService.IsExist(calendarId, taskId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusNotFound, ErrTaskNotFount)
	}

	var body dto.TaskRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return t.taskService.UpdateTask(taskId, &body)

}

func NewTaskHandler(taskService service.TaskService, calService service.CalendarService) TaskHandler {
	return &taskHandler{
		taskService: taskService,
		calService:  calService,
	}
}
