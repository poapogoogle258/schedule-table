package repository

import (
	"errors"
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(userId string) *dao.Users
	FineUserEmail(email string) *dao.Users
	UpdateOneById(userId string, column string, value any) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u *UserRepositoryImpl) FindUser(userId string) *dao.Users {
	var user *dao.Users
	if err := u.db.Find(&user, "id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return user
}

func (u *UserRepositoryImpl) FineUserEmail(email string) *dao.Users {
	var user *dao.Users

	if err := u.db.Model(&dao.Users{}).Find(&user, "email = ?", email).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return user
}

func (u *UserRepositoryImpl) UpdateOneById(userId string, column string, value any) error {
	return u.db.Model(&dao.Users{}).Where("id = ?", userId).Update(column, value).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}
