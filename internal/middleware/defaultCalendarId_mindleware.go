package middleware

import (
	"schedule_table/internal/handler"

	"github.com/gin-gonic/gin"
)

func DefaultCalendarId(handlerCal handler.CalendarsHandler) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Param("calendarId") == "default" {
			userId := c.Keys["token_userId"].(string)
			calendarId := handlerCal.GetDefaultCalendarId(userId)
			c.AddParam("calendarId", calendarId)
		}

		c.Next()

	}
}
