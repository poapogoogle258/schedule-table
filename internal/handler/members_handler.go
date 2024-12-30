package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/constant"
	"schedule_table/internal/interface/request"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/gin-gonic/gin"
)

type MemberHandler interface {
	GetMembers(c *gin.Context)
	GetMemberId(c *gin.Context)
	CreateNewMember(c *gin.Context)
}

type memberHandler struct {
	memberRepo repository.MembersRepository
}

func (mh *memberHandler) GetMembers(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if members, err := mh.memberRepo.GetMembers(c.Param("calendarId")); err != nil {
		panic(err)
	} else {
		c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, members))
	}

}

func (mh *memberHandler) GetMemberId(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if member, err := mh.memberRepo.GetMemberId(c.Param("memberId")); err != nil {
		panic(err)
	} else {
		c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, member))
	}
}

func (mh *memberHandler) CreateNewMember(c *gin.Context) {
	defer pkg.PanicHandler(c)

	var req request.CreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		panic(err)
	}
	if err := req.Validate(); err != nil {
		panic(err)
	}

	req.File.Filename = fmt.Sprintf(`%s(%s).%s`, req.Name, req.NickName, util.GetExpressionFile(req.File.Filename))
	c.SaveUploadedFile(req.File, "../../upload/images/profile/"+req.File.Filename)

	if newMember, err := mh.memberRepo.CreateNewMember(c.Param("calendarId"), &req); err != nil {
		panic(err)
	} else {
		c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, newMember))
	}
}

func NewMemberHandler(memberRepo repository.MembersRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
	}
}
