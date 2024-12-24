package router

import (
	"fmt"
	"net/http"
	"schedule_table/internal/handler"
	"schedule_table/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Calendar     handler.CalendarsHandler
	GenerateTask handler.GenerateTaskHandler
	Auth         handler.AuthHandler
}

func NewRouter(handlers *Handlers) *gin.Engine {

	router := gin.New()

	// CorsConfig := cors.DefaultConfig()
	// CorsConfig.AllowAllOrigins = true
	// CorsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	// CorsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	// CorsConfig.ExposeHeaders = []string{"Content-Length"}
	// CorsConfig.AllowCredentials = true
	// CorsConfig.MaxAge = 60 * 60
	router.Use(cors.Default())

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(CustomRecovery))

	auth := router.Group("/auth")

	{
		auth.POST("/login", handlers.Auth.Login)
		auth.GET("/validate", handlers.Auth.ValidateToken)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthorizeJWT(handlers.Auth))

	{
		calendar := api.Group("/calendars/:calendarId")
		calendar.Use(middleware.DefaultCalendarId(handlers.Calendar))

		calendar.GET("/", handlers.Calendar.GetMyCalendar)
		calendar.POST("/generate", handlers.GenerateTask.GenerateTasks)

	}

	return router
}

func CustomRecovery(c *gin.Context, err any) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"message":    fmt.Sprintf("StatusInternalServerError: %s", err),
	})
}
