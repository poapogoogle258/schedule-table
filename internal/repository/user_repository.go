package repository

import (
	"errors"
	"os"
	"schedule_table/internal/model/dao"
	"schedule_table/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindOne(userId string) (*dao.Users, error)
	FindOneByEmail(email string) (*dao.Users, error)
	UpdateOne(userId string, column string, value any) error
	Profile(userId string) (*dao.Users, error)
	IsUniqueEmail(email string) bool
	Register(insert *dao.Users) error
	CreateCalendarDefault(userId uuid.UUID) (*dao.Calendars, error)
	GetTokenUser(userId string) (string, error)
}

type userRepository struct {
	db *gorm.DB
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrUserNotFound   = errors.New("not found user")
)

func (userRepo *userRepository) GetTokenUser(userId string) (string, error) {
	var user *dao.Users

	result := userRepo.db.Model(&dao.Users{}).Select("id", "token").First(&user, userId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", ErrUserNotFound
		} else {
			return "", result.Error
		}
	}

	return user.Token, nil

}

func (userRepo *userRepository) CreateCalendarDefault(userId uuid.UUID) (*dao.Calendars, error) {
	calendar := &dao.Calendars{
		Id:       uuid.New(),
		Name:     "default",
		UserId:   userId,
		ImageURL: os.Getenv("HOST") + "/upload/default-member-profile.jpeg",
	}

	if err := userRepo.db.Create(&calendar).Error; err != nil {
		return nil, err
	}

	return calendar, nil

}

func (userRepo *userRepository) Register(insert *dao.Users) error {

	// hast password
	insert.Password = util.HashPassword(insert.Password)

	// set userId
	insert.Id = uuid.New()

	return userRepo.db.Create(insert).Error
}

func (userRepo *userRepository) IsUniqueEmail(email string) bool {
	var count int64
	if err := userRepo.db.Model(&dao.Users{}).Limit(1).Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	return count == 0
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
