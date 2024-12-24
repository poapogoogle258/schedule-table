package handler

import (
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
	GetDefaultCalendarId(userId string) string
}

type calendarsHandler struct {
	calRepo repository.CalendarRepository
}

func (s *calendarsHandler) GetMyCalendar(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if userId, check := c.Keys["token_userId"]; check {
		calendar := s.calRepo.FindByOwnerId(userId.(string))
		c.JSON(200, pkg.BuildResponse(constant.Success, calendar))
	} else {
		pkg.PanicException(constant.DataNotFound)
	}
}

func (s *calendarsHandler) GetDefaultCalendarId(userId string) string {
	return s.calRepo.GetDefaultCalendarId(userId)
}

func NewCalendarsHandler(calRepo repository.CalendarRepository) CalendarsHandler {
	return &calendarsHandler{
		calRepo: calRepo,
	}
}
