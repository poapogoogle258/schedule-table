package handler

import (
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context) (*[]dto.ResponseCalendar, error)
}

type calendarsHandler struct {
	calRepo repository.CalendarRepository
}

func (calHandler *calendarsHandler) GetMyCalendar(c *gin.Context) (*[]dto.ResponseCalendar, error) {
	userId := c.GetString("requestAuthUserId")

	return calHandler.calRepo.FindByOwnerId(userId)
}

func NewCalendarsHandler(calRepo repository.CalendarRepository) CalendarsHandler {
	return &calendarsHandler{
		calRepo: calRepo,
	}
}
