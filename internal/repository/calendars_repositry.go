package repository

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"

	"time"

	"gorm.io/gorm"
)

type CalendarRepository interface {
	FindMembersOfCalendarId(calendarId string) (*[]dao.Members, error)
	FindLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) (*[]dao.Leaves, error)
	IsOwnerOfCalendar(userId string, calendarId string) bool
	FindByOwnerId(ownerId string) (*[]dto.ResponseCalendar, error)
	IsExits(calendarId string) bool
	FindOneWithAssociation(calendarId string, start time.Time, end time.Time) (*dao.Calendars, error)
}

type calendarRepository struct {
	db *gorm.DB
}

func (calRepo *calendarRepository) FindByOwnerId(ownerId string) (*[]dto.ResponseCalendar, error) {
	var calendars *[]dto.ResponseCalendar

	if err := calRepo.db.Model(&dao.Calendars{}).Find(&calendars, "user_id = ?", ownerId).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}

func (calRepo *calendarRepository) IsOwnerOfCalendar(userId string, calendarId string) bool {
	var count int64
	calRepo.db.Model(&dao.Calendars{}).Where("id = ? AND user_id = ?", calendarId, userId).Count(&count)

	return count > 0
}

func (calRepo *calendarRepository) FindLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) (*[]dao.Leaves, error) {
	var leaves *[]dao.Leaves
	result := calRepo.db.Model(&dao.Leaves{}).
		Where("calendar_id = ?", calendarId).
		Where("(start BETWEEN ? AND ?) OR (end BETWEEN ? AND ?)", start, end, start, end).
		Find(&leaves)

	if result.Error != nil {
		return nil, result.Error
	}

	return leaves, nil
}

func (calRepo *calendarRepository) FindMembersOfCalendarId(calendarId string) (*[]dao.Members, error) {
	var members *[]dao.Members
	if err := calRepo.db.Preload("Leaves").Find(&members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (calRepo *calendarRepository) IsExits(calendarId string) bool {
	var count int64
	if err := calRepo.db.Model(&dao.Calendars{}).Where("id = ?", calendarId).Count(&count).Error; err != nil {
		panic(err)
	}

	return count > 0

}

func (calRepo *calendarRepository) FindOneWithAssociation(calendarId string, start time.Time, end time.Time) (*dao.Calendars, error) {
	calendar := &dao.Calendars{}

	if err := calRepo.db.
		Preload("Members").
		Preload("Leaves", func(db *gorm.DB) *gorm.DB {
			return db.Where("leaves.date BETWEEN ? AND ?", start, end).Order("leaves.date ASC")
		}).
		Preload("Schedules.Responsibles", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Person").Order("responsibles.queue ASC")
		}).
		First(&calendar, "calendars.id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return calendar, nil

}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db: db,
	}
}
