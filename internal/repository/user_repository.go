package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindOne(userId string) (*dao.Users, error)
	FindOneByEmail(email string) (*dao.Users, error)
	UpdateOne(userId string, column string, value any) error
	Profile(userId string) (*dao.Users, error)
}

type userRepository struct {
	db *gorm.DB
}

func (userRepo *userRepository) FindOne(userId string) (*dao.Users, error) {
	var user *dao.Users
	if err := userRepo.db.Find(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *userRepository) FindOneByEmail(email string) (*dao.Users, error) {
	var user *dao.Users
	if err := userRepo.db.Model(&dao.Users{}).Find(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *userRepository) UpdateOne(userId string, column string, value any) error {
	return userRepo.db.Model(&dao.Users{}).Where("id = ?", userId).Update(column, value).Error
}

func (repo *userRepository) Profile(userId string) (*dao.Users, error) {
	var user *dao.Users
	if err := repo.db.Preload("Calendar").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
