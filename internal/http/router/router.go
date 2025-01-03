package router

import (
	"fmt"
	"net/http"
	"schedule_table/internal/handler"
	"schedule_table/internal/http/middleware"
	"schedule_table/internal/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Calendar handler.CalendarsHandler
	Auth     handler.AuthHandler
	Member   handler.MemberHandler
	Schedule handler.ScheduleHandler
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
		calendar := api.Group("/calendars")

		// calendar
		calendar.GET("/", pkg.BuildGetController(handlers.Calendar.GetMyCalendar))

		// members
		calendar.GET("/:calendarId/members", pkg.BuildGetController(handlers.Member.GetMembers))
		calendar.POST("/:calendarId/members", pkg.BuildPostController(handlers.Member.CreateNewMember))
		calendar.GET("/:calendarId/members/:memberId", pkg.BuildGetController(handlers.Member.GetMemberId))
		calendar.PATCH("/:calendarId/members/:memberId", pkg.BuildPatchController(handlers.Member.EditMember))
		calendar.DELETE("/:calendarId/members/:memberId", pkg.BuildDeleteController(handlers.Member.DeleteMemberId))

		// schedule
		calendar.GET("/:calendarId/schedules", pkg.BuildGetController(handlers.Schedule.GetSchedules))
		calendar.GET("/:calendarId/schedules/:scheduleId", pkg.BuildGetController(handlers.Schedule.GetScheduleId))
		calendar.POST("/:calendarId/schedules", pkg.BuildPostController(handlers.Schedule.CreateNewSchedule))

	}

	return router
}

func CustomRecovery(c *gin.Context, err any) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"message":    fmt.Sprintf("StatusInternalServerError: %s", err),
	})
}
