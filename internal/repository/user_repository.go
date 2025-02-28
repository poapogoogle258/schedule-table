package repository

import (
	"errors"
	"os"
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[*dao.User]
	InjectionTx(tx *gorm.DB) UserRepository
}

type userRepository struct {
	Repository[*dao.User]
	db *gorm.DB
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrUserNotFound   = errors.New("not found user")
)

func (userRepo *userRepository) InjectionTx(tx *gorm.DB) UserRepository {

	return &userRepository{
		db:         tx,
		Repository: NewRepository[*dao.User](tx),
	}
}

func (userRepo *userRepository) GetTokenUser(userId string) (string, error) {
	var user *dao.User

	result := userRepo.db.Model(&dao.User{}).Select("id", "token").First(&user, "id = ?", userId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", ErrUserNotFound
		} else {
			return "", result.Error
		}
	}

	return user.Token, nil

}

func (userRepo *userRepository) CreateCalendarDefault(userId uuid.UUID) (*dao.Calendar, error) {
	calendar := &dao.Calendar{
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

func (userRepo *userRepository) Create(insert *dao.User) error {
	return userRepo.db.Create(insert).Error
}

func (userRepo *userRepository) IsUniqueEmail(email string) bool {
	var count int64
	if err := userRepo.db.Model(&dao.User{}).Limit(1).Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	return count == 0
}

func (userRepo *userRepository) FindOneByEmail(email string) (*dao.User, error) {
	var user *dao.User
	if err := userRepo.db.Model(&dao.User{}).Find(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) GetProfile(userId string) (*dao.User, error) {
	var user *dao.User
	if err := repo.db.Preload("Calendar").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {

	return &userRepository{
		db:         db,
		Repository: NewRepository[*dao.User](db),
	}

}
