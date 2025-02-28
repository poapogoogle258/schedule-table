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
	Repository[*dao.Calendar]
	InjectionTx(tx *gorm.DB) CalendarRepository
}

type calendarRepository struct {
	Repository[*dao.Calendar]
	db *gorm.DB
}

func (repo *calendarRepository) InjectionTx(tx *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db:         tx,
		Repository: NewRepository[*dao.Calendar](tx),
	}
}

func (repo *calendarRepository) GetTasks(calendarId string, start time.Time, end time.Time) ([]*dao.Task, error) {
	var tasks []*dao.Task
	if err := repo.db.Preload("Description").Find(&tasks, "calendar_id = ? AND start BETWEEN ? AND ?", calendarId, start, end).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *calendarRepository) GetSchedules(calendarId string) ([]*dao.Schedule, error) {
	var schedules []*dao.Schedule
	if err := repo.db.Model(&dao.Schedule{}).Preload("Responsibles").Find(&schedules, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (calRepo *calendarRepository) GetDefaultCalendarUser(userId string) (string, error) {

	var calendar *dao.Calendar
	if err := calRepo.db.Select("user_id", "id").First(&calendar, "user_id = ?", userId).Error; err != nil {
		return "", err
	}

	return calendar.Id.String(), nil
}

func (calRepo *calendarRepository) FindByOwnerId(ownerId string) ([]*dto.ResponseCalendar, error) {
	var calendars []*dto.ResponseCalendar

	if err := calRepo.db.Model(&dao.Calendar{}).Find(&calendars, "user_id = ?", ownerId).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}

func (calRepo *calendarRepository) IsOwnerOfCalendar(userId string, calendarId string) bool {
	var count int64
	calRepo.db.Model(&dao.Calendar{}).Where("id = ? AND user_id = ?", calendarId, userId).Count(&count)

	return count > 0
}

func (calRepo *calendarRepository) FindLeavesOfCalendarId(calendarId string, start time.Time, end time.Time) ([]*dao.Leave, error) {
	var leaves []*dao.Leave
	result := calRepo.db.Model(&dao.Leave{}).
		Where("calendar_id = ? AND date BETWEEN ? AND ?", calendarId, start, end).
		Find(&leaves)

	if result.Error != nil {
		return nil, result.Error
	}

	return leaves, nil
}

func (calRepo *calendarRepository) FindMembersOfCalendarId(calendarId string) ([]*dao.Employee, error) {
	var members []*dao.Employee
	if err := calRepo.db.Preload("Leaves").Find(&members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (calRepo *calendarRepository) CheckExist(calendarId string) error {
	var count int64
	if err := calRepo.db.Model(&dao.Calendar{}).Where("id = ?", calendarId).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		return ErrCalendarNotFount
	} else {
		return nil
	}

}

func (calRepo *calendarRepository) FindOneWithAssociation(calendarId string, start time.Time, end time.Time) (*dao.Calendar, error) {
	calendar := &dao.Calendar{}

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
	var calendars []*dao.Calendar
	if err := repo.db.Select("id", "schedule_changed_at", "generate_task_updated_at").Find(&calendars, "schedule_changed_at > generate_task_updated_at").Error; err != nil {
		return nil, err
	}

	result := make([]string, len(calendars))
	for i := range calendars {
		result[i] = calendars[i].Id.String()
	}

	return result, nil

}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db:         db,
		Repository: NewRepository[*dao.Calendar](db),
	}
}
