package handler

import (
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	userId, _ := uuid.Parse(c.Param("userId"))

	calendar := s.calRepo.FindByOwnerId(userId)
	if calendar == nil {
		c.JSON(200, calendar)
	} else {
		c.JSON(200, calendar)
	}

}

func NewCalendarsHandler(calRepo repository.CalendarRepository) CalendarsHandler {
	return &CalendarsHandlerImpl{
		calRepo: calRepo,
	}
}
