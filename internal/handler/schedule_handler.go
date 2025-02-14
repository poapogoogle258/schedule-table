package handler

import (
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type ScheduleHandler interface {
	GetSchedules(c *gin.Context) ([]*dto.ResponseSchedule, error)
	GetScheduleId(c *gin.Context) (*dto.ResponseSchedule, error)
	CreateNewSchedule(c *gin.Context) (*dto.ResponseSchedule, error)
	UpdateSchedule(c *gin.Context) (*dto.ResponseSchedule, error)
	DeleteSchedule(c *gin.Context) error
	GetResponsible(c *gin.Context) ([]*dto.ResponseMember, error)
}

type scheduleHandler struct {
	scheduleRepo repository.ScheduleRepository
}

func (scheHandler *scheduleHandler) GetResponsible(c *gin.Context) ([]*dto.ResponseMember, error) {
	scheduleId := c.Param("scheduleId")

	members, errGetMembersResponsible := scheHandler.scheduleRepo.GetMembersResponsible(scheduleId)
	if errGetMembersResponsible != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errGetMembersResponsible)
	}

	resp := []*dto.ResponseMember{}
	copier.Copy(&resp, members)

	return resp, nil
}

func (scheHandler *scheduleHandler) GetSchedules(c *gin.Context) ([]*dto.ResponseSchedule, error) {
	calendarId := c.Param("calendarId")

	schedules, errGetSchedulesCalendar := scheHandler.scheduleRepo.GetSchedulesCalendar(calendarId)
	if errGetSchedulesCalendar != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errGetSchedulesCalendar)
	}

	response := []*dto.ResponseSchedule{}
	copier.Copy(&response, schedules)

	return response, nil
}

func (scheHandler *scheduleHandler) GetScheduleId(c *gin.Context) (*dto.ResponseSchedule, error) {
	calendarId := c.Param("calendarId")
	scheduleId := c.Param("scheduleId")

	schedule, errGetScheduleCalendarId := scheHandler.scheduleRepo.GetScheduleCalendarId(calendarId, scheduleId)
	if errGetScheduleCalendarId != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errGetScheduleCalendarId)

	}

	resp := dto.ResponseSchedule{}
	copier.Copy(&resp, schedule)

	return &resp, nil

}

func (scheHandler *scheduleHandler) CreateNewSchedule(c *gin.Context) (*dto.ResponseSchedule, error) {
	calendarId := c.Param("calendarId")

	var req dto.RequestSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	// TODO : Validate request

	insert := dao.Schedules{}
	insert.CalendarId = uuid.MustParse(calendarId)
	if err := copier.Copy(&insert, &req); err != nil {
		return nil, err
	}

	newSchedule, errCreateNewSchedule := scheHandler.scheduleRepo.CreateNewSchedule(&insert)
	if errCreateNewSchedule != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errCreateNewSchedule)
	}

	resp := dto.ResponseSchedule{}
	copier.Copy(&resp, newSchedule)

	return &resp, nil

}

func (scheHandler *scheduleHandler) UpdateSchedule(c *gin.Context) (*dto.ResponseSchedule, error) {
	scheduleId := c.Param("scheduleId")

	var req dto.RequestSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	// TODO : Validate request

	insert := dao.Schedules{}
	copier.Copy(&insert, req)

	schedule, errUpdateSchedule := scheHandler.scheduleRepo.UpdateSchedule(scheduleId, &insert)
	if errUpdateSchedule != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errUpdateSchedule)
	}

	response := dto.ResponseSchedule{}
	copier.Copy(&response, schedule)

	return &response, nil

}

func (scheHandler *scheduleHandler) DeleteSchedule(c *gin.Context) error {
	scheduleId := c.Param("scheduleId")

	return scheHandler.scheduleRepo.Delete(scheduleId)

}

func NewScheduleHandler(scheduleRepo repository.ScheduleRepository) ScheduleHandler {
	return &scheduleHandler{
		scheduleRepo: scheduleRepo,
	}
}
