package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUser() ([]dao.Users, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u UserRepositoryImpl) FindAllUser() ([]dao.Users, error) {

}

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}
