package handler

import (
	"schedule_table/internal/model/dto"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context) ([]*dto.CalendarInfo, error)
}

type calendarsHandler struct {
	calService service.CalendarService
}

func (handler *calendarsHandler) GetMyCalendar(c *gin.Context) ([]*dto.CalendarInfo, error) {
	userId := c.GetString("authUserId")

	return handler.calService.GetCalendarsOfUser(userId)
}

func NewCalendarsHandler(calendarService service.CalendarService) CalendarsHandler {
	return &calendarsHandler{
		calService: calendarService,
	}
}
