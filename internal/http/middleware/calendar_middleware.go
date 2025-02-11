package middleware

import (
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type ICalendarMiddleware interface {
	CheckExist() func(c *gin.Context)
}

type CalendarMiddleware struct {
	CalRepo repository.CalendarRepository
}

func (calMiddle *CalendarMiddleware) CheckExist() func(c *gin.Context) {

	return func(c *gin.Context) {
		calendarId := c.Param("calendarId")
		if !calMiddle.CalRepo.IsExist(calendarId) {
			c.JSON(http.StatusNotFound, pkg.BuildWithoutResponse(http.StatusNotFound, repository.ErrCalendarNotFount.Error()))

			c.Abort()
		}
	}
}

func NewCalendarMiddleware(calRepo repository.CalendarRepository) ICalendarMiddleware {
	return &CalendarMiddleware{
		CalRepo: calRepo,
	}
}
