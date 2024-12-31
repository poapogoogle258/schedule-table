package repository

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	GetMyCalendars(ownerId string) (*[]dto.ResponseCalendar, error)
	IsOwnerCalendar(userId string, calendarId string) bool
	GetLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) *[]dao.Leaves
	GetMembersOfCalendarId(calendarId string) *[]dao.Members
}

type CalendarRepositoryImpl struct {
	db *gorm.DB
}

func (s *CalendarRepositoryImpl) GetMyCalendars(ownerId string) (*[]dto.ResponseCalendar, error) {
	var calendars *[]dto.ResponseCalendar
	user_uuid, _ := uuid.Parse(ownerId)

	if err := s.db.Model(&dao.Calendars{}).Find(&calendars, "user_id = ?", user_uuid).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}

func (s *CalendarRepositoryImpl) IsOwnerCalendar(userId string, calendarId string) bool {
	var calendar *dao.Calendars
	s.db.Select("id").Find(&calendar, "id = ? AND user_id = ?", calendarId, userId)

	return calendar != nil

}

func (s *CalendarRepositoryImpl) GetLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) *[]dao.Leaves {
	var leaves *[]dao.Leaves
	s.db.Find(&leaves)

	return leaves
}

func (s *CalendarRepositoryImpl) GetMembersOfCalendarId(calendarId string) *[]dao.Members {
	var members *[]dao.Members
	s.db.Preload("Leaves").Find(&members, "calendar_id = ?", calendarId)

	return members
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &CalendarRepositoryImpl{
		db: db,
	}
}
