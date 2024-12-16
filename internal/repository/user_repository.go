package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(userId string) *dao.Users
	FineUserEmail(email string) *dao.Users
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u *UserRepositoryImpl) FindUser(userId string) *dao.Users {
	var user *dao.Users

	u.db.Model(&dao.Users{}).Find(&user, "id = ?", userId)

	return user
}

func (u *UserRepositoryImpl) FineUserEmail(email string) *dao.Users {
	var user *dao.Users

	u.db.Model(&dao.Users{}).Find(&user, "email = ?", email)

	return user
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}
