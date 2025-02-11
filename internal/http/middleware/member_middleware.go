package middleware

import (
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type IMemberMiddleware interface {
	CheckExist() func(c *gin.Context)
}

type MemberMiddleware struct {
	MemberRepo repository.MembersRepository
}

func (memberMiddle *MemberMiddleware) CheckExist() func(c *gin.Context) {

	return func(c *gin.Context) {
		scheduleId := c.Param("memberId")
		if !memberMiddle.MemberRepo.IsExist(scheduleId) {
			c.JSON(http.StatusNotFound, pkg.BuildWithoutResponse(http.StatusNotFound, repository.ErrMemberNotFount.Error()))

			c.Abort()
		}

		c.Next()
	}
}

func NewMemberMiddleware(memberRepo repository.MembersRepository) IMemberMiddleware {
	return &MemberMiddleware{
		MemberRepo: memberRepo,
	}
}
