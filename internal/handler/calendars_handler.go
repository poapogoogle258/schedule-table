package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/teambition/rrule-go"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
	GenerateTasks(c *gin.Context)
}

type CalendarsHandlerImpl struct {
	calRepo      repository.CalendarRepository
	scheduleRepo repository.SchedulesRepository
	recurService service.RecurrenceService
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
	// leaves := s.calRepo.GetLeavesOfCalendarId(calendarId, &start, &end)

	for _, schedule := range *schedules {
		config := s.recurService.NewScheduleRuleConfig(&schedule, &start, &end)
		r, _ := rrule.NewRRule(*config)

		fmt.Println(r.All())

	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, schedules))

}

func NewCalendarsHandler(
	calRepo repository.CalendarRepository,
	scheduleRepo repository.SchedulesRepository,
	recurService service.RecurrenceService) CalendarsHandler {
	return &CalendarsHandlerImpl{
		calRepo:      calRepo,
		scheduleRepo: scheduleRepo,
		recurService: recurService,
	}
}
