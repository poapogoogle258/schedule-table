package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	GetSchedulesCalendar(calendarId string) (*[]dao.Schedules, error)
	GetScheduleCalendarId(calendarId string, scheduleId string) (*dao.Schedules, error)
	CreateNewSchedule(insert *dao.Schedules) (*dao.Schedules, error)
}

type scheduleRepository struct {
	db *gorm.DB
}

func (scheduleRepo *scheduleRepository) GetSchedulesCalendar(calendarId string) (*[]dao.Schedules, error) {
	var schedules *[]dao.Schedules

	// TODO: Order by Quese
	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}

	if err := scheduleRepo.db.Model(&dao.Schedules{}).Preload("Responsibles.Person", scheduleRepo.db.Select(selectedField)).Find(&schedules, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (scheduleRepo *scheduleRepository) GetScheduleCalendarId(calendarId string, scheduleId string) (*dao.Schedules, error) {
	var schedule *dao.Schedules
	if err := scheduleRepo.db.Model(&dao.Schedules{}).Preload("Responsibles.Person").First(&schedule, "id = ? AND calendar_id = ?", scheduleId, calendarId).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (scheduleRepo *scheduleRepository) CreateNewSchedule(insert *dao.Schedules) (*dao.Schedules, error) {
	if err := scheduleRepo.db.Create(&insert).Error; err != nil {
		return nil, err
	}

	var schedule *dao.Schedules
	if err := scheduleRepo.db.Model(&dao.Schedules{}).Preload("Responsibles.Person").First(&schedule, "id = ?", insert.Id).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}
