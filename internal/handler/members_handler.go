package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type QueryStringGetMembers struct {
	Page  int  `form:"page"`
	Limit int  `form:"limit"`
	All   bool `form:"all"`
}

func (mh *memberHandler) GetMembers(c *gin.Context) (*dto.ResponseMembersTable, error) {
	var query QueryStringGetMembers
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	calendarId := c.Param("calendarId")
	if err := mh.calRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	if query.All {
		result, err := mh.memberRepo.Find(map[string]interface{}{
			"calendar_id": calendarId,
		})
		if err != nil {
			return nil, err
		}

		response := util.Convert[[]dto.ResponseMember](&result)

		return &dto.ResponseMembersTable{
			Data: response,
			Pagination: &dto.Pagination{
				Total: int64(len(*response)),
			},
		}, nil

	} else {
		result, err := mh.memberRepo.FindWithOffsetAndLimit((query.Page-1)*query.Limit, query.Limit, map[string]interface{}{
			"calendar_id": calendarId,
		})
		if err != nil {
			return nil, err
		}

		response := util.Convert[[]dto.ResponseMember](&result)
		totalMember := mh.memberRepo.Count(calendarId)

		return &dto.ResponseMembersTable{
			Data: response,
			Pagination: &dto.Pagination{
				CurrentPage: query.Page,
				Limit:       query.Limit,
				Total:       totalMember,
			},
		}, nil
	}

}

func (mh *memberHandler) GetMemberId(c *gin.Context) (*dto.ResponseMember, error) {

	calendarId := c.Param("calendarId")
	if err := mh.calRepo.CheckExist(calendarId); err != nil {
		return nil, err
	} else {
		fmt.Println("err : ", err)
	}

	memberId := c.Param("memberId")
	if err := mh.memberRepo.CheckExist(memberId); err != nil {
		return nil, err
	}

	result, err := mh.memberRepo.FindOne(map[string]interface{}{
		"id":          memberId,
		"calendar_id": calendarId,
	})

	if err != nil {
		return nil, err
	}

	response := util.Convert[dto.ResponseMember](&result)

	return response, nil
}

func (mh *memberHandler) CreateNewMember(c *gin.Context) (*dto.ResponseMember, error) {
	var req dto.RequestCreateNewMember
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}
	if err := req.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	calendarId := c.Param("calendarId")
	if err := mh.calRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	insert := util.Convert[dao.Members](&req)
	insert.CalendarId = uuid.MustParse(calendarId)

	if err := mh.memberRepo.Create(insert); err != nil {
		return nil, err
	}

	response := util.Convert[dto.ResponseMember](&insert)

	return response, nil
}

func (mh *memberHandler) EditMember(c *gin.Context) (*dto.ResponseMember, error) {

	var req dto.RequestCreateNewMember
	if err := c.ShouldBind(&req); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}
	if err := req.Validate(); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	calendarId := c.Param("calendarId")
	if err := mh.calRepo.CheckExist(calendarId); err != nil {
		return nil, err
	}

	memberId := c.Param("memberId")
	if err := mh.memberRepo.CheckExist(memberId); err != nil {
		return nil, err
	}

	data := util.Convert[dao.Members](&req)

	if result, err := mh.memberRepo.UpdatesAndFindOne(memberId, calendarId, data); err != nil {
		return nil, err
	} else {
		response := util.Convert[dto.ResponseMember](&result)
		return response, nil
	}

}

func (mh *memberHandler) DeleteMemberId(c *gin.Context) error {
	calendarId := c.Param("calendarId")
	if err := mh.calRepo.CheckExist(calendarId); err != nil {
		return err
	}

	memberId := c.Param("memberId")
	if err := mh.memberRepo.CheckExist(memberId); err != nil {
		return err
	}

	return mh.memberRepo.DeleteOne(memberId, calendarId)
}

func NewMemberHandler(memberRepo repository.MembersRepository, calRepo repository.CalendarRepository) MemberHandler {
	return &memberHandler{
		memberRepo: memberRepo,
		calRepo:    calRepo,
	}
}
