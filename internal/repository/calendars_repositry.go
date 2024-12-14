package repository

import (
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	db *gorm.DB
}

func (c *CalendarRepository) FindByOwnerId(ownerId uuid.UUID) *dao.Calendars {

	var calendar *dao.Calendars

	user_uuid, _ := uuid.Parse("238e921d-8362-4232-847c-bf747465cdaf")

	c.db.Preload("Leaves").Preload("Schedules").Preload("Tasks").Preload("Members").Find(&calendar, "user_id = ?", user_uuid)

	return calendar

}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{
		db: db,
	}
}
