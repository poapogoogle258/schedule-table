package router

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"schedule_table/internal/handler"
	"schedule_table/internal/http/middleware"
	"schedule_table/internal/pkg"
	"schedule_table/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth           handler.AuthHandler
	AuthJWTMiddle  middleware.IAuthorizeJWTMiddleware
	Calendar       handler.CalendarsHandler
	CalendarMiddle middleware.ICalendarMiddleware
	Member         handler.MemberHandler
	MemberMiddle   middleware.IMemberMiddleware
	Schedule       handler.ScheduleHandler
	ScheduleMiddle middleware.IScheduleMiddleware
	Task           handler.TasksHandler
	TaskMiddle     middleware.ITaskMiddleware
}

func NewRouter(handlers *Handlers) *gin.Engine {

	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma", "x-requested-with"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(CustomRecovery))

	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Auth.Login)
		auth.POST("/signup", handlers.Auth.SignUp)
		auth.GET("/profile", handlers.Auth.Profile)
	}

	api := router.Group("/api")
	api.Use(handlers.AuthJWTMiddle.Authorize())

	{
		calendar := api.Group("/calendars/:calendarId")

		{
			calendar.Use(handlers.CalendarMiddle.CheckExist(), handlers.CalendarMiddle.IsOwner())

			// members
			calendar.GET("/members", pkg.BuildGetController(handlers.Member.GetMembers))
			calendar.POST("/members", pkg.BuildPostController(handlers.Member.CreateNewMember))
			calendar.GET("/members/:memberId", handlers.MemberMiddle.CheckExist(), pkg.BuildGetController(handlers.Member.GetMemberId))
			calendar.PATCH("/members/:memberId", handlers.MemberMiddle.CheckExist(), pkg.BuildPatchController(handlers.Member.EditMember))
			calendar.DELETE("/members/:memberId", handlers.MemberMiddle.CheckExist(), pkg.BuildDeleteController(handlers.Member.DeleteMemberId))

			// schedule
			calendar.GET("/schedules", pkg.BuildGetController(handlers.Schedule.GetSchedules))
			calendar.POST("/schedules", pkg.BuildPostController(handlers.Schedule.CreateNewSchedule))
			calendar.GET("/schedules/:scheduleId", handlers.ScheduleMiddle.CheckExist(), pkg.BuildGetController(handlers.Schedule.GetScheduleId))
			calendar.PATCH("/schedules/:scheduleId", handlers.ScheduleMiddle.CheckExist(), pkg.BuildPatchController(handlers.Schedule.UpdateSchedule))
			calendar.DELETE("/schedules/:scheduleId", handlers.ScheduleMiddle.CheckExist(), pkg.BuildDeleteController(handlers.Schedule.DeleteSchedule))

			// get members responsible
			calendar.GET("/schedules/:scheduleId/responsible", handlers.ScheduleMiddle.CheckExist(), pkg.BuildGetController(handlers.Schedule.GetResponsible))

			// task
			calendar.GET("/tasks", pkg.BuildGetController(handlers.Task.GetTasks))
			calendar.PATCH("/tasks/:taskId", handlers.TaskMiddle.CheckExist(), pkg.BuildPatchController(handlers.Task.EditTask))

			// Not this phaser 1 implement on phaser 2 or 3 ?
			// leave
			// calendar.GET("/leaves", pkg.BuildGetController(handlers.Leave.GetLeave))
			// calendar.POST("/leaves", pkg.BuildPostController(handlers.Leave.CreateNewLeave))
			// calendar.DELETE("/leaves/:leaveId", pkg.BuildDeleteController(handlers.Leave.Delete))
		}

	}

	type Form struct {
		File *multipart.FileHeader `form:"avatar" binding:"required"`
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
			Url:      os.Getenv("HOST") + "/upload/" + form.File.Filename,
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
