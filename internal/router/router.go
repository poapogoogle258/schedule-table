package router

import (
	"fmt"
	"net/http"
	"schedule_table/internal/handler"
	"schedule_table/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Calendar handler.CalendarsHandler
	Auth     handler.AuthHandler
}

func NewRouter(handlers *Handlers) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(CustomRecovery))

	auth := router.Group("/auth")

	{
		auth.POST("/login", handlers.Auth.Login)
		auth.GET("/validate", handlers.Auth.ValidateToken)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthorizeJWT())

	{
		user := api.Group("/calendars")
		user.GET("/default", handlers.Calendar.GetMyCalendar)
	}

	return router
}

func CustomRecovery(c *gin.Context, err any) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"message":    fmt.Sprintf("StatusInternalServerError: %s", err),
	})
}
