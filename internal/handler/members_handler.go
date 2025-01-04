package handler

import (
	"errors"
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
	response := &[]dto.ResponseMember{}

	result, err := mh.memberRepo.FindByCalendarId(c.Param("calendarId"))
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil
}

func (mh *memberHandler) GetMemberId(c *gin.Context) (*dto.ResponseMember, error) {

	response := &dto.ResponseMember{}

	result, err := mh.memberRepo.FindOne(c.Param("memberId"))
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil
}

func (mh *memberHandler) CreateNewMember(c *gin.Context) (*dto.ResponseMember, error) {
	var req dto.RequestCreateNewMember

	// TODO : must send calendarId in from request
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

	result, err := mh.memberRepo.Create(insert)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseMember{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}

	return response, nil
}

func (mh *memberHandler) EditMember(c *gin.Context) (*dto.ResponseMember, error) {
	memberId := c.Param("memberId")
	calendarId := c.Param("calendarId")

	var req dto.RequestCreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if !mh.memberRepo.IsMemberExits(memberId) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("not fount member id in calendar"))
	}

	if req.File != nil {
		req.File.Filename = fmt.Sprintf(`%v.%s`, time.Now().UnixMicro(), util.GetExpressionFile(req.File.Filename))
		c.SaveUploadedFile(req.File, "../../upload/images/profile/"+req.File.Filename)
	}

	insertData := &dao.Members{}
	copier.Copy(insertData, &req)
	insertData.CalendarId = uuid.Must(uuid.Parse(calendarId))

	result, err := mh.memberRepo.UpdateOne(memberId, insertData)
	if err != nil {
		return nil, err
	}

	response := &dto.ResponseMember{}
	if err := copier.Copy(&response, &result); err != nil {
		return nil, err
	}
	return response, nil
}

func (mh *memberHandler) DeleteMemberId(c *gin.Context) error {
	memberId := c.Param("memberId")

	if !mh.memberRepo.IsMemberExits(memberId) {
		return pkg.NewErrorWithStatusCode(http.StatusBadRequest, errors.New("not fount member id in calendar"))
	}

	return mh.memberRepo.DeleteOne(memberId)
}

func NewMemberHandler(memberRepo repository.MembersRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
	}
}
