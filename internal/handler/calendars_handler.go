package handler

import (
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
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
	defer pkg.PanicHandler(c)

	if userId, check := c.Keys["token_userId"]; check {
		calendar := s.calRepo.FindByOwnerId(userId.(string))
		c.JSON(200, pkg.BuildResponse(constant.Success, calendar))
	} else {
		pkg.PanicException(constant.DataNotFound)
	}

}

func NewCalendarsHandler(calRepo *repository.CalendarRepositoryImpl) *CalendarsHandlerImpl {
	return &CalendarsHandlerImpl{
		calRepo: calRepo,
	}
}
