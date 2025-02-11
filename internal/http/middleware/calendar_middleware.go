package middleware

import (
	"errors"
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type ICalendarMiddleware interface {
	CheckExist() func(c *gin.Context)
	IsOwner() func(c *gin.Context)
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

		c.Next()
	}
}

func (calMiddle *CalendarMiddleware) IsOwner() func(c *gin.Context) {
	return func(c *gin.Context) {
		calendarId := c.Param("calendarId")
		userId := c.GetString("authUserId")

		if defaultCalendarId, err := calMiddle.CalRepo.GetDefaultCalendarUser(userId); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.BuildWithoutResponse(http.StatusInternalServerError, err.Error()))
			c.Abort()

			return
		} else {
			if defaultCalendarId != calendarId {
				c.JSON(http.StatusForbidden, pkg.BuildWithoutResponse(http.StatusForbidden, ErrForbidden.Error()))
				c.Abort()

				return
			}

			c.Next()
		}

	}
}

func NewCalendarMiddleware(calRepo repository.CalendarRepository) ICalendarMiddleware {
	return &CalendarMiddleware{
		CalRepo: calRepo,
	}
}
