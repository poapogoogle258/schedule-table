package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type MemberHandler interface {
	GetMembers(c *gin.Context) (*[]dto.ResponseMember, error)
	GetMemberId(c *gin.Context) (*dto.ResponseMember, error)
	CreateNewMember(c *gin.Context) (*dto.ResponseMember, error)
	EditMember(c *gin.Context) (*dto.ResponseMember, error)
	DeleteMemberId(c *gin.Context) error
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

	if req.File != nil {
		req.File.Filename = fmt.Sprintf(`%v.%s`, time.Now().UnixMicro(), util.GetExpressionFile(req.File.Filename))
		c.SaveUploadedFile(req.File, "../../upload/images/profile/"+req.File.Filename)
	}

	insert := &dao.Members{}

	copier.Copy(insert, &req)
	insert.CalendarId = uuid.Must(uuid.Parse(c.Param("calendarId")))

	return mh.memberRepo.CreateNewMember(insert)
}

func (mh *memberHandler) EditMember(c *gin.Context) (*dto.ResponseMember, error) {
	memberId := c.Param("memberId")
	calendarId := c.Param("calendarId")

	var req dto.RequestCreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.ErrorWithStatusCode{Code: http.StatusBadRequest, Err: err}
	}

	if err := req.Validate(); err != nil {
		return nil, pkg.ErrorWithStatusCode{Code: http.StatusBadRequest, Err: err}
	}

	if !mh.memberRepo.ExistMemberId(calendarId, memberId) {
		return nil, pkg.ErrorWithStatusCode{Code: http.StatusBadRequest, Err: fmt.Errorf("EditMember: not fount member id in calendar")}
	}

	if req.File != nil {
		req.File.Filename = fmt.Sprintf(`%v.%s`, time.Now().UnixMicro(), util.GetExpressionFile(req.File.Filename))
		c.SaveUploadedFile(req.File, "../../upload/images/profile/"+req.File.Filename)
	}

	insertData := &dao.Members{}
	copier.Copy(insertData, &req)

	insertData.CalendarId = uuid.Must(uuid.Parse(calendarId))

	return mh.memberRepo.EditMember(memberId, insertData)

}

func (mh *memberHandler) DeleteMemberId(c *gin.Context) error {
	memberId := c.Param("memberId")
	calendarId := c.Param("calendarId")

	if !mh.memberRepo.ExistMemberId(calendarId, memberId) {
		return pkg.ErrorWithStatusCode{Code: http.StatusBadRequest, Err: fmt.Errorf("DeleteMemberId: not fount member id in calendar")}
	}

	return mh.memberRepo.DeleteMemberId(calendarId, memberId)

}

func NewMemberHandler(memberRepo repository.MembersRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
	}
}
