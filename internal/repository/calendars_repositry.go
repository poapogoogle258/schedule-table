package repository

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	FindByOwnerId(ownerId string) *dao.Calendars
	IsOwnerCalendar(userId string, calendarId string) bool
	GetDefaultCalendarId(userId string) string
	GetLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) *[]dao.Leaves
	GetMembersOfCalendarId(calendarId string) *[]dao.Members
}

type CalendarRepositoryImpl struct {
	db *gorm.DB
}

func (s *CalendarRepositoryImpl) FindByOwnerId(ownerId string) *dao.Calendars {

	var calendar *dao.Calendars

	user_uuid, _ := uuid.Parse(ownerId)

	s.db.Preload("Leaves").Preload("Schedules").Preload("Tasks").Preload("Members").Find(&calendar, "user_id = ?", user_uuid)

	return calendar

}

func (s *CalendarRepositoryImpl) IsOwnerCalendar(userId string, calendarId string) bool {
	var calendar *dao.Calendars
	s.db.Select("id").Find(&calendar, "id = ? AND user_id = ?", calendarId, userId)

	return calendar != nil

}

func (s *CalendarRepositoryImpl) GetDefaultCalendarId(userId string) string {
	var calendar *dao.Calendars
	s.db.Select("id").Find(&calendar, "user_id = ? AND name = ?", "default", userId)

	return calendar.Id.String()

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
