package router

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"schedule_table/internal/handler"
	"schedule_table/internal/http/middleware"
	"schedule_table/internal/pkg"
	"schedule_table/util"
	"time"

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
		calendar.PATCH("/:calendarId/schedules/:scheduleId", pkg.BuildPatchController(handlers.Schedule.UpdateSchedule))
		calendar.DELETE("/:calendarId/schedules/:scheduleId", pkg.BuildDeleteController(handlers.Schedule.DeleteSchedule))

	}

	type Form struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}
	router.POST("/upload", func(c *gin.Context) {
		var form Form
		if err := c.ShouldBind(&form); err != nil {
			panic(err)
		}

		form.File.Filename = fmt.Sprintf(`%v.%s`, time.Now().UnixMicro(), util.GetExpressionFile(form.File.Filename))
		c.SaveUploadedFile(form.File, "../../upload/public/"+form.File.Filename)

		c.JSON(http.StatusOK, pkg.BuildResponse(http.StatusOK, struct {
			Filename string `json:"filename"`
			Url      string `json:"url"`
		}{
			Filename: form.File.Filename,
			Url:      "/upload/" + form.File.Filename,
		}))
		c.Abort()
	})
	router.Static("/upload", "../../upload/public")

	return router
}

func CustomRecovery(c *gin.Context, err any) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"message":    fmt.Sprintf("StatusInternalServerError: %s", err),
	})
}
