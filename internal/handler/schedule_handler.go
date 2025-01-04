package handler

import (
	"errors"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type ScheduleHandler interface {
	GetSchedules(c *gin.Context) (*[]dto.ResponseSchedule, error)
	GetScheduleId(c *gin.Context) (*dto.ResponseSchedule, error)
	CreateNewSchedule(c *gin.Context) (*dto.ResponseSchedule, error)
	UpdateSchedule(c *gin.Context) (*dto.ResponseSchedule, error)
	DeleteSchedule(c *gin.Context) error
}

type scheduleHandler struct {
	scheduleRepo repository.ScheduleRepository
}

func (scheHandler *scheduleHandler) GetSchedules(c *gin.Context) (*[]dto.ResponseSchedule, error) {
	calendarId := c.Param("calendarId")

	result, err := scheHandler.scheduleRepo.GetSchedulesCalendar(calendarId)
	if err != nil {
		return nil, err
	}

	response := &[]dto.ResponseSchedule{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil
}

func (scheHandler *scheduleHandler) GetScheduleId(c *gin.Context) (*dto.ResponseSchedule, error) {
	calendarId := c.Param("calendarId")
	scheduleId := c.Param("scheduleId")

	result, err := scheHandler.scheduleRepo.GetScheduleCalendarId(calendarId, scheduleId)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseSchedule{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil

}

func (scheHandler *scheduleHandler) CreateNewSchedule(c *gin.Context) (*dto.ResponseSchedule, error) {
	// calendarId := c.Param("calendarId")
	var req *dto.RequestSchedule
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, errors.New("bad request"))
	}

	// TODO : Validate request

	insert := &dao.Schedules{}
	if err := copier.Copy(&insert, &req); err != nil {
		return nil, err
	}

	result, err := scheHandler.scheduleRepo.CreateNewSchedule(insert)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseSchedule{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil

}

func (scheHandler *scheduleHandler) UpdateSchedule(c *gin.Context) (*dto.ResponseSchedule, error) {
	scheduleId := c.Param("scheduleId")

	var req *dto.RequestSchedule
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(400, err)
	}

	// TODO : Validate request
	// TODO : CheckExit schedule id

	insert := &dao.Schedules{}
	if err := copier.Copy(&insert, &req); err != nil {
		return nil, err
	}

	result, err := scheHandler.scheduleRepo.UpdateSchedule(scheduleId, insert)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseSchedule{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil

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
