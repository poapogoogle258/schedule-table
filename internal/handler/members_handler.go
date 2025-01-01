package handler

import (
	"fmt"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/gin-gonic/gin"
)

type MemberHandler interface {
	GetMembers(c *gin.Context) (*[]dto.ResponseMember, error)
	GetMemberId(c *gin.Context) (*dto.ResponseMember, error)
	CreateNewMember(c *gin.Context) (*dto.ResponseMember, error)
}

type memberHandler struct {
	memberRepo repository.MembersRepository
}

func (mh *memberHandler) GetMembers(c *gin.Context) (*[]dto.ResponseMember, error) {

	return mh.memberRepo.GetMembers(c.Param("calendarId"))
}

func (mh *memberHandler) GetMemberId(c *gin.Context) (*dto.ResponseMember, error) {

	return mh.memberRepo.GetMemberId(c.Param("memberId"))
}

func (mh *memberHandler) CreateNewMember(c *gin.Context) (*dto.ResponseMember, error) {
	var req dto.RequestCreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		panic(err)
	}
	if err := req.Validate(); err != nil {
		panic(err)
	}

	req.File.Filename = fmt.Sprintf(`%s(%s).%s`, req.Name, req.NickName, util.GetExpressionFile(req.File.Filename))
	c.SaveUploadedFile(req.File, "../../upload/images/profile/"+req.File.Filename)

	return mh.memberRepo.CreateNewMember(c.Param("calendarId"), &req)
}

func NewMemberHandler(memberRepo repository.MembersRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
	}
}
