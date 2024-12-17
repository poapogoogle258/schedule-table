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
	s.db.Select("id").Find(&calendar, "name = ? AND user_id = ?", "default", userId)

	return calendar.Id.String()

}

func (s *CalendarRepositoryImpl) GetLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) *[]dao.Leaves {
	var leaves *[]dao.Leaves
	s.db.Find(&leaves)

	return leaves
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &CalendarRepositoryImpl{
		db: db,
	}
}
