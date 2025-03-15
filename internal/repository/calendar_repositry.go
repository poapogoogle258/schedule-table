package repository

import (
	"errors"
	"schedule_table/internal/model/dao"
	"time"

	"gorm.io/gorm"
)

var (
	ErrCalendarNotFount = errors.New("not fount calendar")
)

type CalendarRepository interface {
	Repository[*dao.Calendar]
	InjectionTx(tx *gorm.DB) CalendarRepository
	FindCalendarIdWithFullAggregate(calendarId string, start time.Time, end time.Time) (*dao.Calendar, error)
}

type calendarRepository struct {
	Repository[*dao.Calendar]
	db *gorm.DB
}

func (repo *calendarRepository) FindCalendarIdWithFullAggregate(calendarId string, start time.Time, end time.Time) (*dao.Calendar, error) {
	var calendar *dao.Calendar

	result := repo.db.
		Preload("Leaves", func(db *gorm.DB) *gorm.DB {
			return db.Where("date between ? and ?", start, end)
		}).
		Preload("Schedules", func(db *gorm.DB) *gorm.DB {
			return db.Preload("EmployeeQueue", func(db *gorm.DB) *gorm.DB {
				return db.Order("employee_queues.queue ASC")
			})
		}).
		Preload("Employees").
		First(&calendar, "id = ?", calendarId)

	if result.Error != nil {
		return nil, result.Error
	}

	return calendar, nil

}

func (repo *calendarRepository) InjectionTx(tx *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db:         tx,
		Repository: NewRepository[*dao.Calendar](tx),
	}
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db:         db,
		Repository: NewRepository[*dao.Calendar](db),
	}
}
