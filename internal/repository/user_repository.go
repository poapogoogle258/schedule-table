package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(userId string) *dao.Users
	FineUserEmail(email string) *dao.Users
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) FindUser(userId string) *dao.Users {
	var user *dao.Users

	u.db.Model(&dao.Users{}).Find(&user, "id = ?", userId)

	return user
}

func (u *userRepository) FineUserEmail(email string) *dao.Users {
	var user *dao.Users

	u.db.Model(&dao.Users{}).Find(&user, "email = ?", email)

	return user
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
