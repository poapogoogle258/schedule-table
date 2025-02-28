package service

import (
	"errors"
	"regexp"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserService interface {
	GetProfileByUserId(userId string) (*dao.User, error)
	ValidateNewUser(name string, email string, password string, image string) error
	RegisterUser(request *dto.SignUpRequest) (*dto.SignUpResponse, error)
	Authentication(email string, password string) (*dao.User, error)
	UpdateToken(userId, token string) error
	GetProfile(userId string) (*dto.UserProfile, error)
}

type userService struct {
	transaction  repository.Transaction
	userRepo     repository.UserRepository
	calendarRepo repository.CalendarRepository
}

func (s *userService) GetProfileByUserId(userId string) (*dao.User, error) {
	return s.userRepo.FindOne("id = ?", userId)
}

func (s *userService) GetProfileByEmail(email string) (*dao.User, error) {
	return s.userRepo.FindOne("email = ?", email)
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (s *userService) ValidateNewUser(name string, email string, password string, image string) error {
	if len(name) <= 0 {
		return errors.New("validateNewUser: length name mush more than 0")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("validateNewUser: email invalid")
	}

	if s.userRepo.Count("email = ?", email) > 0 {
		return errors.New("validateNewUser: email already exists")
	}

	if len(password) < 8 {
		return errors.New("validateNewUser: length password mush more than 8")
	}

	return nil
}

func (s *userService) RegisterUser(request *dto.SignUpRequest) (*dto.SignUpResponse, error) {

	tx := s.transaction.Begin()
	defer tx.Rollback()

	userRepo := s.userRepo.InjectionTx(tx)
	calendarRepo := s.calendarRepo.InjectionTx(tx)

	user := &dao.User{
		Id:          uuid.New(),
		Name:        request.Name,
		ImageURL:    request.ImageUrl,
		Email:       request.Email,
		Password:    util.HashPassword(request.Password),
		Description: request.Description,
	}

	defaultCalendarNewUser := &dao.Calendar{
		Id:       user.Id,
		UserId:   user.Id,
		Name:     "default",
		ImageURL: constant.DEFAULT_IMAGE_PROFILE,
	}

	if err := userRepo.Create(user); err != nil {
		return nil, err
	}

	if err := calendarRepo.Create(defaultCalendarNewUser); err != nil {
		return nil, err
	}

	result := dto.SignUpResponse{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, err
	}

	tx.Commit()
	return &result, nil

}

var ErrAuthEmailOrPasswordInvalid = errors.New("email or password invalid")

func (s *userService) Authentication(email string, password string) (*dao.User, error) {
	user, err := s.userRepo.FindOne("email = ?", email)
	if err != nil {
		return nil, ErrAuthEmailOrPasswordInvalid
	}

	if !util.VerifyPassword(password, user.Password) {
		return nil, ErrAuthEmailOrPasswordInvalid
	}

	return user, nil

}

func (s *userService) GetProfile(userId string) (*dto.UserProfile, error) {
	user, err := s.userRepo.FindOneWithAggregate([]string{"Calendar"}, "id = ?", userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		ImageURL:    user.ImageURL,
		Description: user.Description,
		CalendarId:  user.Calendar.Id.String(),
	}, nil

}

func (s *userService) UpdateToken(userId, token string) error {
	return s.userRepo.UpdateColumn(userId, "token", token)
}

func NewUserService(userRepo repository.UserRepository, calendarRepo repository.CalendarRepository, transaction repository.Transaction) UserService {
	return &userService{
		userRepo:     userRepo,
		calendarRepo: calendarRepo,
		transaction:  transaction,
	}
}
