package handler

import (
	"errors"
	"net/http/httptest"
	"schedule_table/internal/handler"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	mocks "schedule_table/test/mock"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMemberId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {

		mockMemberRepo := new(mocks.MembersRepository)
		mockCalRepo := new(mocks.CalendarRepository)
		memberHandler := handler.NewMemberHandler(mockMemberRepo, mockCalRepo)

		calendarId := "c080f838-694f-4e9a-b6c6-1e552decf7a0"
		memberId := "1f65b672-dbc0-4048-b64e-05cbb9c712d8"
		expect := &dto.ResponseMember{
			Id:          "1f65b672-dbc0-4048-b64e-05cbb9c712d8",
			ImageURL:    "https://robohash.org/quiablanditiislaborum.png?size=200x200&set=set1",
			Name:        "Berkly Maun",
			Nickname:    "Berkly",
			Color:       "#e5559a",
			Description: "Unsp injury at C4 level of cervical spinal cord, init encntr",
			Position:    "Infinite Therapies of Sarasota, Inc.",
			Email:       "bmaun1@gmpg.org",
			Telephone:   "410 886 3030",
		}
		mockCalRepo.On("CheckExist", calendarId).Return(nil)
		mockMemberRepo.On("CheckExist", memberId).Return(nil)
		mockMemberRepo.On("FindOne", map[string]interface{}{
			"id":          memberId,
			"calendar_id": calendarId,
		}).Return(&dao.Members{
			Id:          uuid.MustParse("1f65b672-dbc0-4048-b64e-05cbb9c712d8"),
			ImageURL:    "https://robohash.org/quiablanditiislaborum.png?size=200x200&set=set1",
			Name:        "Berkly Maun",
			Nickname:    "Berkly",
			Color:       "#e5559a",
			Description: "Unsp injury at C4 level of cervical spinal cord, init encntr",
			Position:    "Infinite Therapies of Sarasota, Inc.",
			Email:       "bmaun1@gmpg.org",
			Telephone:   "410 886 3030",
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "calendarId", Value: calendarId},
			{Key: "memberId", Value: memberId},
		}

		resp, err := memberHandler.GetMemberId(c)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expect, resp)
	})

	t.Run("Calendar Not Exist", func(t *testing.T) {

		mockMemberRepo := new(mocks.MembersRepository)
		mockCalRepo := new(mocks.CalendarRepository)
		memberHandler := handler.NewMemberHandler(mockMemberRepo, mockCalRepo)

		calendarId := "c080f838-694f-4e9a-b6c6-1e552decf7a0"
		memberId := "1f65b672-dbc0-4048-b64e-05cbb9c712d8"

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "calendarId", Value: calendarId},
			{Key: "memberId", Value: memberId},
		}

		mockCalRepo.On("CheckExist", calendarId).Return(repository.ErrCalendarNotFount)

		resp, err := memberHandler.GetMemberId(c)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, err, repository.ErrCalendarNotFount)

	})

	t.Run("Member Not Exist", func(t *testing.T) {

		mockMemberRepo := new(mocks.MembersRepository)
		mockCalRepo := new(mocks.CalendarRepository)
		memberHandler := handler.NewMemberHandler(mockMemberRepo, mockCalRepo)

		calendarId := "c080f838-694f-4e9a-b6c6-1e552decf7a0"
		memberId := "1f65b672-dbc0-4048-b64e-05cbb9c712d8"

		mockCalRepo.On("CheckExist", calendarId).Return(nil)
		mockMemberRepo.On("CheckExist", memberId).Return(repository.ErrMemberNotFount)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "calendarId", Value: calendarId},
			{Key: "memberId", Value: memberId},
		}

		resp, err := memberHandler.GetMemberId(c)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, err, repository.ErrMemberNotFount)
	})

	t.Run("FindOne Error", func(t *testing.T) {
		mockMemberRepo := new(mocks.MembersRepository)
		mockCalRepo := new(mocks.CalendarRepository)
		memberHandler := handler.NewMemberHandler(mockMemberRepo, mockCalRepo)

		calendarId := "c080f838-694f-4e9a-b6c6-1e552decf7a0"
		memberId := "1f65b672-dbc0-4048-b64e-05cbb9c712d8"

		mockCalRepo.On("CheckExist", calendarId).Return(nil)
		mockMemberRepo.On("CheckExist", memberId).Return(nil)
		mockMemberRepo.On("FindOne", mock.Anything).Return(&dao.Members{}, errors.New("find error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "calendarId", Value: calendarId},
			{Key: "memberId", Value: memberId},
		}

		resp, err := memberHandler.GetMemberId(c)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
