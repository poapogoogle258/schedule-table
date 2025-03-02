package handler

import (
	"schedule_table/internal/model/dto"
	"schedule_table/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type LeaveHandler interface {
	FindAllLeavesOfCalendar(c *gin.Context) ([]*dto.TaskInfo, error)
}

type leaveHandler struct {
	leaveService service.LeaveService
}

type findTasksOfRangeQuery struct {
	start time.Time `form:"start" binding:"required"`
	end   time.Time `form:"end" binding:"required"`
}

func (h *leaveHandler) FindAllLeavesOfCalendar(c *gin.Context) ([]*dto.TaskInfo, error) {
	calendarId := c.Param("calendarId")
	var query findTasksOfRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	return h.leaveService.FindAllLeavesOfCalendar(calendarId, query.start, query.end)
}

func NewLeaveHandler(leaveService service.LeaveService) LeaveHandler {
	return &leaveHandler{
		leaveService: leaveService,
	}
}
