package repository

import (
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	FindByOwnerId(ownerId string) *dao.Calendars
}

type CalendarRepositoryImpl struct {
	db *gorm.DB
}

func (c *CalendarRepositoryImpl) FindByOwnerId(ownerId string) *dao.Calendars {

	var calendar *dao.Calendars

	user_uuid, _ := uuid.Parse(ownerId)

	c.db.Preload("Leaves").Preload("Schedules").Preload("Tasks").Preload("Members").Find(&calendar, "user_id = ?", user_uuid)

	return calendar

}

func NewCalendarRepository(db *gorm.DB) *CalendarRepositoryImpl {
	return &CalendarRepositoryImpl{
		db: db,
	}
}
