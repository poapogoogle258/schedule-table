package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler interface {
	GetAllSchedules(c *gin.Context) ([]*dto.ScheduleInfo, error)
	GetSchedule(c *gin.Context) (*dto.ScheduleInfo, error)
	CreateSchedule(c *gin.Context) (*dto.ScheduleInfo, error)
	UpdateSchedule(c *gin.Context) (*dto.ScheduleInfo, error)
	DeleteSchedule(c *gin.Context) error
}

type scheduleHandler struct {
	transaction     repository.Transaction
	employeeService service.EmployeeService
	scheduleService service.ScheduleService
}

var ErrScheduleNotFound = errors.New("schedule not found")

func (handler *scheduleHandler) GetAllSchedules(c *gin.Context) ([]*dto.ScheduleInfo, error) {
	calendarId := c.Param("calendarId")

	return handler.scheduleService.GetAllSchedules(calendarId)
}

func (handler *scheduleHandler) GetSchedule(c *gin.Context) (*dto.ScheduleInfo, error) {
	scheduleId := c.Param("scheduleId")

	return handler.scheduleService.GetSchedule(scheduleId)
}

var ErrEmployeeNotFound = errors.New("employee not found")
var ErrMasterScheduleNotFound = errors.New("master schedule not found")

func (handler *scheduleHandler) CreateSchedule(c *gin.Context) (*dto.ScheduleInfo, error) {
	calendarId := c.Param("calendarId")

	var body dto.ScheduleInfoRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if len(body.MasterScheduleId) != 0 && !handler.scheduleService.IsExist(body.MasterScheduleId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrMasterScheduleNotFound)
	}

	employeeIds := make([]string, len(body.Employees))
	for i := range employeeIds {
		employeeIds[i] = body.Employees[i].Id
	}
	if !handler.employeeService.EmployeesIsExistCalendar(calendarId, employeeIds) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrEmployeeNotFound)
	}

	return handler.scheduleService.CreateSchedule(calendarId, &body)

	// start transaction and start generate tasks to calendar view
}

func (handler *scheduleHandler) UpdateSchedule(c *gin.Context) (*dto.ScheduleInfo, error) {

	scheduleId := c.Param("scheduleId")
	if !handler.scheduleService.IsExist(scheduleId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrScheduleNotFound)
	}

	var body dto.ScheduleInfoRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return handler.scheduleService.UpdateSchedule(scheduleId, &body)

	// start transaction and start generate tasks to calendar view

}

func (handler *scheduleHandler) DeleteSchedule(c *gin.Context) error {

	scheduleId := c.Param("scheduleId")
	if !handler.scheduleService.IsExist(scheduleId) {
		return pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrScheduleNotFound)
	}

	// start transaction and delete task and relation before delete schedule

	return handler.scheduleService.DeleteSchedule(scheduleId)

}

func NewScheduleHandler(transaction repository.Transaction, scheduleService service.ScheduleService, employeeService service.EmployeeService) ScheduleHandler {
	return &scheduleHandler{
		transaction:     transaction,
		scheduleService: scheduleService,
		employeeService: employeeService,
	}
}
