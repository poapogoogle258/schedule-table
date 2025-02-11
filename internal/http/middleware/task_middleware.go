package middleware

import (
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type ITaskMiddleware interface {
	CheckExist() func(c *gin.Context)
}

type TaskMiddleware struct {
	TaskRepo repository.ITaskRepository
}

func (taskMiddle *TaskMiddleware) CheckExist() func(c *gin.Context) {

	return func(c *gin.Context) {
		calendarId := c.Param("taskId")
		if !taskMiddle.TaskRepo.IsExist(calendarId) {
			c.JSON(http.StatusNotFound, pkg.BuildWithoutResponse(http.StatusNotFound, repository.ErrCalendarNotFount.Error()))

			c.Abort()
		}

		c.Next()
	}
}

func NewTaskMiddleware(taskRepo repository.ITaskRepository) ITaskMiddleware {
	return &TaskMiddleware{
		TaskRepo: taskRepo,
	}
}
