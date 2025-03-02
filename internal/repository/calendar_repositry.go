package repository

import (
	"errors"
	"schedule_table/internal/model/dao"

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

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepository{
		db:         db,
		Repository: NewRepository[*dao.Calendar](db),
	}
}
