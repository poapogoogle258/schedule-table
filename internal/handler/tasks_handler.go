package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type TasksHandler interface {
	GetTasks(c *gin.Context) (*dao.Calendars, error)
}

type tasksHandler struct {
	CalRepo repository.CalendarRepository
}

type queryStringGetTasks struct {
	Start  string `form:"start" binding:"required"`
	End    string `form:"end" binding:"required"`
	Action string `form:"action" binding:"required"`
}

func (taskHandler *tasksHandler) GetTasks(c *gin.Context) (*dao.Calendars, error) {
	var query queryStringGetTasks
	if err := c.BindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("query string not validate"))
	}

	calendarId := c.Param("calendarId")
	if !taskHandler.CalRepo.IsExits(calendarId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("not found calendar id"))
	}

	// start := util.Must(time.Parse(time.RFC3339, query.Start))
	// end := util.Must(time.Parse(time.RFC3339, query.End))

	calendar, errFindCalendar := taskHandler.CalRepo.FindOneWithAssociation(calendarId)
	if errFindCalendar != nil {
		return nil, errFindCalendar
	}

	return calendar, nil

}

func NewTasksHandler(calRepo repository.CalendarRepository) TasksHandler {
	return &tasksHandler{
		CalRepo: calRepo,
	}
}
