package handler

import (
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type MemberHandler interface {
}

type MemberHandlerImpl struct {
	userRepo repository.UserRepository
}

func (mh *MemberHandlerImpl) GetMembers(c *gin.Context) {

}

func (mh *MemberHandlerImpl) GetMemberId(c *gin.Context) {

}

func (mh *MemberHandlerImpl) CreateNewMember(c *gin.Context) {

}

func (mh *MemberHandlerImpl) EditMemberInfo(c *gin.Context) {

}

func (mh *MemberHandlerImpl) DeleteMember(c *gin.Context) {

}

func NewMemberHandler(userRepo repository.UserRepository) MemberHandler {
	return &MemberHandlerImpl{
		userRepo: userRepo,
	}
}
