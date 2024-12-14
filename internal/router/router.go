package router

import (
	"schedule_table/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Calendar handler.CalendarsHandler
}

func NewRouter(handlers *Handlers) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")

	{
		user := api.Group("/calendars")
		user.GET("/default", handlers.Calendar.GetMyCalendar)
	}

	return router
}
