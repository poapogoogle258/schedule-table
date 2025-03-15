package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrLeaveNotFound = errors.New("leave not found")
)

type LeaveHandler interface {
	FindAllLeavesOfCalendar(c *gin.Context) ([]*dto.LeaveInfo, error)
	FindLeaveIdOfCalendar(c *gin.Context) (*dto.LeaveInfo, error)
	CreateLeave(c *gin.Context) (*dto.LeaveInfo, error)
	ChangeLeaveStatus(c *gin.Context) (*dto.LeaveInfo, error)
}

type leaveHandler struct {
	leaveService service.LeaveService
}

type findTasksOfRangeQuery struct {
	start time.Time `form:"start" binding:"required" time_format:"2006-01-02"`
	end   time.Time `form:"end" binding:"required" time_format:"2006-01-02"`
}

func (h *leaveHandler) FindAllLeavesOfCalendar(c *gin.Context) ([]*dto.LeaveInfo, error) {
	calendarId := c.Param("calendarId")
	var query findTasksOfRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	return h.leaveService.FindAllLeavesOfCalendar(calendarId, query.start, query.end)
}

func (h *leaveHandler) FindLeaveIdOfCalendar(c *gin.Context) (*dto.LeaveInfo, error) {
	calendarId := c.Param("calendarId")
	leaveId := c.Param("leaveId")

	if !h.leaveService.IsExistOfCalendar(calendarId, leaveId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusNotFound, ErrLeaveNotFound)
	}

	return h.leaveService.FindOndLeave(leaveId)
}

func (h *leaveHandler) CreateLeave(c *gin.Context) (*dto.LeaveInfo, error) {

	// TO DO: validate date to create and commit
	calendarId := c.Param("calendarId")
	var body = &dto.LeaveRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return h.leaveService.CreateLeave(calendarId, body)

}

func (h *leaveHandler) ChangeLeaveStatus(c *gin.Context) (*dto.LeaveInfo, error) {
	calendarId := c.Param("calendarId")
	leaveId := c.Param("leaveId")
	newStatus := c.Param("newStatus")

	// TO DO: validate date to create and commit

	if !h.leaveService.IsExistOfCalendar(calendarId, leaveId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusNotFound, ErrLeaveNotFound)
	}

	if err := h.leaveService.ChangeStatusLeave(leaveId, newStatus); err != nil {
		return nil, err
	}

	return h.leaveService.FindOndLeave(leaveId)

	// TO DO: schedule to check leaveDoc not accept change status to cancel

}

func NewLeaveHandler(leaveService service.LeaveService) LeaveHandler {
	return &leaveHandler{
		leaveService: leaveService,
	}
}
