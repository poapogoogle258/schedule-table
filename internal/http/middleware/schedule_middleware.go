package middleware

import (
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type IScheduleMiddleware interface {
	CheckExist() func(c *gin.Context)
}

type ScheduleMiddleware struct {
	ScheRepo repository.ScheduleRepository
}

func (scheMiddle *ScheduleMiddleware) CheckExist() func(c *gin.Context) {

	return func(c *gin.Context) {
		scheduleId := c.Param("scheduleId")
		if !scheMiddle.ScheRepo.IsExist(scheduleId) {
			c.JSON(http.StatusNotFound, pkg.BuildWithoutResponse(http.StatusNotFound, repository.ErrScheduleNotExit.Error()))

			c.Abort()
		}
	}
}

func NewScheduleMiddleware(scheRepo repository.ScheduleRepository) IScheduleMiddleware {
	return &ScheduleMiddleware{
		ScheRepo: scheRepo,
	}
}
