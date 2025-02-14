package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var (
	ErrNotFountPagination = errors.New("pagination mush have page and limit in query string or all=true")
)

type MemberHandler interface {
	GetMembers(c *gin.Context) (*dto.ResponseMembersTable, error)
	GetMemberId(c *gin.Context) (*dto.ResponseMember, error)
	CreateNewMember(c *gin.Context) (*dto.ResponseMember, error)
	EditMember(c *gin.Context) (*dto.ResponseMember, error)
	DeleteMemberId(c *gin.Context) error
}

type memberHandler struct {
	memberRepo repository.MembersRepository
	calRepo    repository.CalendarRepository
}

type queryStringGetMembers struct {
	Page  *int  `form:"page"`
	Limit *int  `form:"limit"`
	All   *bool `form:"all"`
}

func (mh *memberHandler) GetMembers(c *gin.Context) (*dto.ResponseMembersTable, error) {
	calendarId := c.Param("calendarId")

	var query queryStringGetMembers
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	resp := []*dto.ResponseMember{}
	pagination := dto.Pagination{}

	if query.All != nil && *query.All {
		members, errFindMember := mh.memberRepo.Find("calendar_id = ?", calendarId)
		if errFindMember != nil {
			return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errFindMember)
		}

		copier.Copy(&resp, members)
		pagination.CurrentPage = 1
		pagination.Limit = len(resp)
		pagination.Total = int64(len(resp))

	} else {
		if query.Page == nil || query.Limit == nil {
			return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, ErrNotFountPagination)
		}

		offset := (*query.Page - 1) * *query.Limit
		limit := *query.Limit
		members, errFindWithOffsetAndLimit := mh.memberRepo.FindWithOffsetAndLimit(offset, limit, "calendar_id = ?", calendarId)
		if errFindWithOffsetAndLimit != nil {
			return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errFindWithOffsetAndLimit)
		}

		copier.Copy(&resp, members)
		pagination.CurrentPage = *query.Page
		pagination.Limit = limit
		pagination.Total = mh.memberRepo.Count(calendarId)
	}

	return &dto.ResponseMembersTable{Data: resp, Pagination: &pagination}, nil

}

func (mh *memberHandler) GetMemberId(c *gin.Context) (*dto.ResponseMember, error) {

	calendarId := c.Param("calendarId")
	memberId := c.Param("memberId")

	member, errFindOne := mh.memberRepo.FindOne("id = ? AND calendar_id = ?", memberId, calendarId)
	if errFindOne != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errFindOne)
	}

	resp := dto.ResponseMember{}
	copier.Copy(&resp, member)

	return &resp, nil
}

func (mh *memberHandler) CreateNewMember(c *gin.Context) (*dto.ResponseMember, error) {

	calendarId := c.Param("calendarId")

	var req dto.RequestCreateNewMember
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := req.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	insert := util.Convert[dao.Members](&req)
	insert.CalendarId = uuid.MustParse(calendarId)

	if err := mh.memberRepo.Create(insert); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	response := util.Convert[dto.ResponseMember](&insert)

	return response, nil
}

func (mh *memberHandler) EditMember(c *gin.Context) (*dto.ResponseMember, error) {

	calendarId := c.Param("calendarId")
	memberId := c.Param("memberId")

	var req dto.RequestCreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}
	if err := req.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	data := util.Convert[dao.Members](&req)

	member, errUpdatesAndFindOne := mh.memberRepo.UpdatesAndFindOne(memberId, calendarId, data)
	if errUpdatesAndFindOne != nil {
		return nil, errUpdatesAndFindOne
	}

	resp := dto.ResponseMember{}
	copier.Copy(&resp, member)

	return &resp, nil

}

func (mh *memberHandler) DeleteMemberId(c *gin.Context) error {

	calendarId := c.Param("calendarId")
	memberId := c.Param("memberId")

	return mh.memberRepo.DeleteOne(memberId, calendarId)
}

func NewMemberHandler(memberRepo repository.MembersRepository, calRepo repository.CalendarRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
		calRepo:    calRepo,
	}
}
