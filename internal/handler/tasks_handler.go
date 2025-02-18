package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/service"
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
	TaskService service.TaskService
}

type queryStringGetTasks struct {
	Start string `form:"start" binding:"required"`
	End   string `form:"end" binding:"required"`
}

func (tasksHandler *tasksHandler) parseAndValidateQuery(c *gin.Context) (start time.Time, end time.Time, err error) {
	var query queryStringGetTasks
	if errShouldBindQuery := c.ShouldBindQuery(&query); errShouldBindQuery != nil {
		err = errShouldBindQuery
		return
	}

	start, errParseStart := time.Parse(time.RFC3339, query.Start)
	if errParseStart != nil {
		err = errParseStart
		return
	}
	end, errParseEnd := time.Parse(time.RFC3339, query.End)
	if errParseEnd != nil {
		err = ErrQueryEndMustFormat
	}

	return
}

func (taskHandler *tasksHandler) GetTasks(c *gin.Context) ([]*dto.ResponseTask, error) {
	calendarId := c.Param("calendarId")

	start, end, err := taskHandler.parseAndValidateQuery(c)
	if err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	tasks, errFindRangeTasksOfCalendarId := taskHandler.TaskService.FindRangeTasksOfCalendarId(calendarId, start, end)
	if errFindRangeTasksOfCalendarId != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errFindRangeTasksOfCalendarId)
	}

	resp := []*dto.ResponseTask{}
	copier.Copy(&resp, tasks)

	return resp, nil

}

var (
	ErrStartDateTimeIsAfterEndDateTime = errors.New("start datetime is after end datetime")
)

type reserveMemberBody struct {
	MemberId *uuid.UUID `json:"member_id" binding:"required"`
	Start    time.Time  `json:"start" binding:"required"`
	End      time.Time  `json:"end" binding:"required"`
	Status   int8       `json:"status" binding:"required"`
}

func (body *reserveMemberBody) Validate() error {
	if body.Start.After(body.End) {
		return ErrStartDateTimeIsAfterEndDateTime
	}

	if _, err := constant.ParseTaskStatus(body.Status); err != nil {
		return err
	}

	return nil
}

func (handler *tasksHandler) EditTask(c *gin.Context) (*dto.ResponseTask, error) {

	taskId := c.Param("taskId")

	body := &reserveMemberBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	if err := body.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	task, errEditTaskId := handler.TaskService.EditTaskId(taskId, body)
	if errEditTaskId != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errEditTaskId)
	}

	resp := dto.ResponseTask{}
	copier.Copy(&resp, task)

	return &resp, nil

}

func NewTasksHandler(taskService service.TaskService) TasksHandler {

	return &tasksHandler{
		TaskService: taskService,
	}
}
