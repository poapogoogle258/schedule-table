package middleware

import (
	"errors"
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type CalendarMiddleware interface {
	CheckExist() func(c *gin.Context)
	IsOwner() func(c *gin.Context)
}

type calendarMiddleware struct {
	calService service.CalendarService
}

func (calMiddle *calendarMiddleware) CheckExist() func(c *gin.Context) {

	return func(c *gin.Context) {
		calendarId := c.Param("calendarId")
		if !calMiddle.calService.IsExist(calendarId) {
			c.JSON(http.StatusNotFound, pkg.BuildWithoutResponse(http.StatusNotFound, repository.ErrCalendarNotFount.Error()))

			c.Abort()
		}

		c.Next()
	}
}

func (calMiddle *calendarMiddleware) IsOwner() func(c *gin.Context) {
	return func(c *gin.Context) {
		calendarId := c.Param("calendarId")
		userId := c.GetString("authUserId")

		if calendarId != userId {
			c.JSON(http.StatusForbidden, pkg.BuildWithoutResponse(http.StatusForbidden, ErrForbidden.Error()))
			c.Abort()

			return
		}

		c.Next()
	}
}

func NewCalendarMiddleware(calService service.CalendarService) CalendarMiddleware {
	return &calendarMiddleware{
		calService: calService,
	}
}
