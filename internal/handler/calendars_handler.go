package handler

import (
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
	GenerateTasks(c *gin.Context)
}

type CalendarsHandlerImpl struct {
	calRepo      repository.CalendarRepository
	scheduleRepo repository.SchedulesRepository
}

func (s *CalendarsHandlerImpl) GetMyCalendar(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if userId, check := c.Keys["token_userId"]; check {
		calendar := s.calRepo.FindByOwnerId(userId.(string))
		c.JSON(200, pkg.BuildResponse(constant.Success, calendar))
	} else {
		pkg.PanicException(constant.DataNotFound)
	}
}

func (s *CalendarsHandlerImpl) GenerateTasks(c *gin.Context) {
	defer pkg.PanicHandler(c)

	userId, validUserId := c.Keys["token_userId"].(string)
	calendarId := c.Param("calendarId")

	if !validUserId {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId != "default" && !s.calRepo.IsOwnerCalendar(userId, calendarId) {
		pkg.PanicException(constant.DataNotFound)
	}

	if calendarId == "default" {
		calendarId = s.calRepo.GetDefaultCalendarId(userId)
	}

	start, end := now.BeginningOfDay(), now.EndOfMonth()

	schedules := s.scheduleRepo.GetScheduleOfCalendar(calendarId, &start, &end)
	leaves := s.calRepo.GetLeavesOfCalendarId(calendarId, &start, &end)

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, schedules))

}

func NewCalendarsHandler(calRepo repository.CalendarRepository, scheduleRepo repository.SchedulesRepository) CalendarsHandler {
	return &CalendarsHandlerImpl{
		calRepo:      calRepo,
		scheduleRepo: scheduleRepo,
	}
}
