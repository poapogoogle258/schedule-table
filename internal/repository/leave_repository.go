package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type LeaveRepository interface {
	Repository[*dao.Leave]
	InjectionTx(tx *gorm.DB) LeaveRepository
}

type leaveRepository struct {
	Repository[*dao.Leave]
	db *gorm.DB
}

func (repo *leaveRepository) InjectionTx(tx *gorm.DB) LeaveRepository {
	return NewLeaveRepository(tx)
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{
		Repository: NewRepository[*dao.Leave](db),
		db:         db,
	}
}
