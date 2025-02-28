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
	CalendarMiddle middleware.CalendarMiddleware
	Employee       handler.EmployeeHandler
	Schedule       handler.ScheduleHandler
}

func NewRouter(handlers *Handlers) *gin.Engine {

	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(CustomRecovery))

	auth := router.Group("/auth")
	{
		auth.GET("/profile", pkg.BuildGetController(handlers.Auth.GetProfile))
		auth.POST("/login", pkg.BuildGetController(handlers.Auth.Login))
		auth.POST("/signup", pkg.BuildPostController(handlers.Auth.SignUp))
	}

	api := router.Group("/api")
	api.Use(handlers.AuthJWTMiddle.Authorize())

	{
		api.GET("/calendars", pkg.BuildGetController(handlers.Calendar.GetMyCalendar))

		calendar := api.Group("/calendars/:calendarId")
		calendar.Use(handlers.CalendarMiddle.CheckExist(), handlers.CalendarMiddle.IsOwner())
		{
			// employee path
			calendar.GET("/employees", pkg.BuildGetController(handlers.Employee.GetEmployees))
			calendar.POST("/employees", pkg.BuildPostController(handlers.Employee.CreateEmployee))
			calendar.GET("/employees/:employeeId", pkg.BuildGetController(handlers.Employee.GetEmployee))
			calendar.PATCH("/employees/:employeeId", pkg.BuildPatchController(handlers.Employee.UpdateEmployee))
			calendar.DELETE("/employees/:employeeId", pkg.BuildDeleteController(handlers.Employee.DeleteEmployee))

			// schedule path
			calendar.GET("/schedules", pkg.BuildGetController(handlers.Schedule.GetAllSchedules))
			calendar.POST("/schedules", pkg.BuildPostController(handlers.Schedule.CreateSchedule))
			calendar.GET("/schedules/:scheduleId", pkg.BuildGetController(handlers.Schedule.GetSchedule))
			calendar.PATCH("/schedules/:scheduleId", pkg.BuildPatchController(handlers.Schedule.UpdateSchedule))
			calendar.DELETE("/schedules/:scheduleId", pkg.BuildDeleteController(handlers.Schedule.DeleteSchedule))

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
