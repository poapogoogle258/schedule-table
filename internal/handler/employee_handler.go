package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/pkg/validator"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"
	"schedule_table/util"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler interface {
	GetEmployees(c *gin.Context) (*EmployeeTableResponse, error)
	GetEmployee(c *gin.Context) (*dto.EmployeeInfo, error)
	CreateEmployee(c *gin.Context) (*dto.EmployeeInfo, error)
	UpdateEmployee(c *gin.Context) (*dto.EmployeeInfo, error)
	DeleteEmployee(c *gin.Context) error
}

type employeeHandler struct {
	transaction     repository.Transaction
	employeeService service.EmployeeService
}

type queryStringGetMembers struct {
	Page  *int `form:"page"`
	Limit *int `form:"limit"`
}

type EmployeeTableResponse struct {
	Data       []*dto.EmployeeInfo `json:"data"`
	Pagination *dto.Pagination     `json:"pagination"`
}

var (
	ErrInvalidPagination = errors.New("invalid query pagination")
	ErrEmployeeNotFount  = errors.New("employee not fount")
)

func (handler *employeeHandler) GetEmployees(c *gin.Context) (*EmployeeTableResponse, error) {
	calendarId := c.Param("calendarId")

	var query queryStringGetMembers
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if query.Limit == nil && query.Page == nil {

		employees := util.Must(handler.employeeService.GetEmployeesOfCalendar(calendarId))
		pagination := dto.Pagination{
			CurrentPage: 1,
			Limit:       len(employees),
			Total:       int64(len(employees)),
		}

		return &EmployeeTableResponse{Data: employees, Pagination: &pagination}, nil

	} else {
		if query.Page == nil || query.Limit == nil {
			return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrInvalidPagination)
		}

		employees := util.Must(handler.employeeService.GetPaginationEmployeesOfCalendarId(calendarId, *query.Page, *query.Limit))
		pagination := dto.Pagination{
			CurrentPage: *query.Page,
			Limit:       *query.Limit,
			Total:       handler.employeeService.CountEmployeesOfCalendarId(calendarId),
		}

		return &EmployeeTableResponse{Data: employees, Pagination: &pagination}, nil
	}

}

func (handler *employeeHandler) GetEmployee(c *gin.Context) (*dto.EmployeeInfo, error) {
	employeeId := c.Param("employeeId")
	if !handler.employeeService.IsExist(employeeId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrEmployeeNotFount)
	}

	return handler.employeeService.GetEmployee(employeeId)

}

func (handler *employeeHandler) CreateEmployee(c *gin.Context) (*dto.EmployeeInfo, error) {
	calendarId := c.Param("calendarId")

	var body dto.EmployeeInfoRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := validator.Validate(body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return handler.employeeService.CreateEmployee(calendarId, &body)
}

func (handler *employeeHandler) UpdateEmployee(c *gin.Context) (*dto.EmployeeInfo, error) {

	employeeId := c.Param("employeeId")
	if !handler.employeeService.IsExist(employeeId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrEmployeeNotFount)
	}

	var body dto.EmployeeInfoRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := validator.Validate(body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return handler.employeeService.UpdateEmployee(employeeId, &body)
}

func (handler *employeeHandler) DeleteEmployee(c *gin.Context) error {

	employeeId := c.Param("employeeId")
	if !handler.employeeService.IsExist(employeeId) {
		return pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrEmployeeNotFount)
	}

	// start transaction
	tx := handler.transaction.Begin()
	defer tx.Rollback()

	employeeService := handler.employeeService.InjectionTx(tx)

	// TO DO: Cancel Tasks And Queue In calendar

	if err := employeeService.DeleteEmployee(employeeId); err != nil {
		return pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	tx.Commit()
	return nil

}

func NewEmployeeHandler(employeeService service.EmployeeService, transaction repository.Transaction) EmployeeHandler {
	return &employeeHandler{
		employeeService: employeeService,
		transaction:     transaction,
	}
}
