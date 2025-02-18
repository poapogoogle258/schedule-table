package repository

import (
	"errors"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"

	"time"

	"gorm.io/gorm"
)

var (
	ErrCalendarNotFount = errors.New("not fount calendar")
)

type CalendarRepository interface {
	FindMembersOfCalendarId(calendarId string) ([]*dao.Members, error)
	FindLeavesOfCalendarId(calendarId string, start *time.Time, end *time.Time) (*[]dao.Leaves, error)
	IsOwnerOfCalendar(userId string, calendarId string) bool
	FindByOwnerId(ownerId string) ([]*dto.ResponseCalendar, error)
	CheckExist(calendarId string) error
	FindOneWithAssociation(calendarId string, start time.Time, end time.Time) (*dao.Calendars, error)
	GetDefaultCalendarUser(userId string) (string, error)
	IsExist(calendarId string) bool
	GetListIdOfScheduleChanged() ([]string, error)
	UpdateScheduleChanged(calendarId string, updateAt time.Time) error
	UpdateTaskGenerated(calendarId string, updateAt time.Time) error
	GetLastTimeGeneratedTask(calendarId string) (*time.Time, error)
}

type calendarRepository struct {
	db *gorm.DB
}

func (calRepo *calendarRepository) GetDefaultCalendarUser(userId string) (string, error) {

	var calendar *dao.Calendars
	if err := calRepo.db.Select("user_id", "id").First(&calendar, "user_id = ?", userId).Error; err != nil {
		return "", err
	}

	return calendar.Id.String(), nil
}

func (calRepo *calendarRepository) IsExist(calendarId string) bool {

	var count int64
	if err := calRepo.db.Model(&dao.Calendars{}).Where("id = ?", calendarId).Limit(1).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (repo *calendarRepository) IsScheduleCalendarChanged(calendarId string) bool {

	var calendar *dao.Calendars
	if err := repo.db.Select("id", "schedule_changed_at", "generate_task_updated_at").First(&calendar, "id = ?", calendarId).Error; err != nil {
		return false
	}

	return calendar.LastTimeGeneratedTask == nil || calendar.LastTimeGeneratedTask.Before(calendar.LastTimeScheduleChanged)
}

func (repo *calendarRepository) UpdateScheduleChanged(calendarId string, updateAt time.Time) error {
	return repo.db.Model(&dao.Calendars{}).Where("id = ?", calendarId).Update("schedule_changed_at", updateAt).Error
}

func (calRepo *calendarRepository) UpdateTaskGenerated(calendarId string, updateAt time.Time) error {
	return calRepo.db.Model(&dao.Calendars{}).Where("id = ?", calendarId).Update("generate_task_updated_at", updateAt).Error
}

func (calRepo *calendarRepository) FindByOwnerId(ownerId string) ([]*dto.ResponseCalendar, error) {
	var calendars []*dto.ResponseCalendar

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

func (calRepo *calendarRepository) FindMembersOfCalendarId(calendarId string) ([]*dao.Members, error) {
	var members []*dao.Members
	if err := calRepo.db.Preload("Leaves").Find(&members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (calRepo *calendarRepository) CheckExist(calendarId string) error {
	var count int64
	if err := calRepo.db.Model(&dao.Calendars{}).Where("id = ?", calendarId).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		return ErrCalendarNotFount
	} else {
		return nil
	}

}

func (calRepo *calendarRepository) FindOneWithAssociation(calendarId string, start time.Time, end time.Time) (*dao.Calendars, error) {
	calendar := &dao.Calendars{}

	if err := calRepo.db.
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Leaves", func(db *gorm.DB) *gorm.DB {
				return db.Where("leaves.date BETWEEN ? AND ?", start, end).Order("leaves.date ASC")
			})
		}).
		Preload("Schedules.Responsibles", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Person").Order("responsibles.queue ASC")
		}).
		First(&calendar, "calendars.id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return calendar, nil

}

func (repo *calendarRepository) GetListIdOfScheduleChanged() ([]string, error) {
	var calendars []*dao.Calendars
	if err := repo.db.Select("id", "schedule_changed_at", "generate_task_updated_at").Find(&calendars, "schedule_changed_at > generate_task_updated_at").Error; err != nil {
		return nil, err
	}

	result := make([]string, len(calendars))
	for i := range calendars {
		result[i] = calendars[i].Id.String()
	}

	return result, nil

}

func (repo *calendarRepository) GetLastTimeGeneratedTask(calendarId string) (*time.Time, error) {
	var calendar *dao.Calendars
	if err := repo.db.Select("id", "generate_task_updated_at").First(&calendar, "id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return calendar.LastTimeGeneratedTask, nil
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db: db,
	}
}
