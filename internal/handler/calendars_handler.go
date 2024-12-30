package handler

import (
	"fmt"
	"schedule_table/internal/constant"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type CalendarsHandler interface {
	GetMyCalendar(c *gin.Context)
}

type calendarsHandler struct {
	calRepo repository.CalendarRepository
}

func (ch *calendarsHandler) GetMyCalendar(c *gin.Context) {
	defer pkg.PanicHandler(c)

	userId := c.GetString("requestAuthUserId")

	if calendars, err := ch.calRepo.GetMyCalendars(userId); err != nil {
		panic(fmt.Errorf("GetMyCalendar: %w", err))
	} else {
		c.JSON(200, pkg.BuildResponse(constant.Success, calendars))
	}
}

func NewCalendarsHandler(calRepo repository.CalendarRepository) CalendarsHandler {
	return &calendarsHandler{
		calRepo: calRepo,
	}
}
