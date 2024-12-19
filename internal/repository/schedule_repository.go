package repository

import (
	"schedule_table/internal/model/dao"
	"time"

	"gorm.io/gorm"
)

type SchedulesRepository interface {
	GetScheduleOfCalendar(calendarId string, start *time.Time, end *time.Time) *[]dao.Schedules
	GetResponsiblePersons(scheduleId string) *[]dao.Responsible
}

type SchedulesRepositoryImpl struct {
	db *gorm.DB
}

func (s *SchedulesRepositoryImpl) GetScheduleOfCalendar(calendarId string, start *time.Time, end *time.Time) *[]dao.Schedules {
	var schedules *[]dao.Schedules

	s.db.Preload("Members_responsible", func(db *gorm.DB) *gorm.DB {
		return db.Select("schedule_id", "member_id", "queue")
	}).
		Find(&schedules, "calendar_id = ?", calendarId)

	return schedules

}

func (s *SchedulesRepositoryImpl) GetResponsiblePersons(scheduleId string) *[]dao.Responsible {
	var responsibles *[]dao.Responsible
	s.db.Order("queue ASC").Find(&responsibles, "schedule_id = ?", scheduleId)

	return responsibles
}

func NewSchedulesRepository(db *gorm.DB) SchedulesRepository {
	return &SchedulesRepositoryImpl{
		db: db,
	}
}
