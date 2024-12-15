package handler

import (
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
}

type CalendarsHandlerImpl struct {
	calRepo repository.CalendarRepository
}

// func (s *CalendarsService) GenerateTasks(start time.Time, end time.Time, calendar dao.Calendars) []*dao.Tasks {
// }

func (s *CalendarsHandlerImpl) GetMyCalendar(c *gin.Context) {

	if userId, check := c.Keys["token_userId"]; check {
		calendar := s.calRepo.FindByOwnerId(userId.(string))
		c.JSON(200, calendar)
	} else {
		c.JSON(200, nil)
	}

}

func NewCalendarsHandler(calRepo repository.CalendarRepository) CalendarsHandler {
	return &CalendarsHandlerImpl{
		calRepo: calRepo,
	}
}
