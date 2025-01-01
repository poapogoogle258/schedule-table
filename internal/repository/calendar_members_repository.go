package repository

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembersRepository interface {
	GetMemberId(memberId string) (*dto.ResponseMember, error)
	GetMembers(calendarId string) (*[]dto.ResponseMember, error)
	CreateNewMember(calendarId string, req *dto.RequestCreateNewMember) (*dto.ResponseMember, error)
}

type membersRepository struct {
	db *gorm.DB
}

func (m *membersRepository) GetMemberId(memberId string) (*dto.ResponseMember, error) {
	var member *dto.ResponseMember

	if err := m.db.Model(&dao.Members{}).First(&member, "id = ?", memberId).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (m *membersRepository) GetMembers(calendarId string) (*[]dto.ResponseMember, error) {
	var Members *[]dto.ResponseMember

	if err := m.db.Model(&dao.Members{}).Find(&Members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return Members, nil
}

func (m *membersRepository) CreateNewMember(calendarId string, req *dto.RequestCreateNewMember) (*dto.ResponseMember, error) {

	newMember := &dao.Members{
		Id:          uuid.New(),
		ImageURL:    req.File.Filename,
		CalendarId:  util.Must(uuid.Parse(calendarId)),
		Name:        req.Name,
		Nickname:    req.NickName,
		Color:       req.Color,
		Description: req.Description,
		Position:    req.Position,
		Email:       req.Email,
		Telephone:   req.Telephone,
	}

	if err := m.db.Model(&dao.Members{}).Create(&newMember).Error; err != nil {
		return nil, err
	}

	var responseMember *dto.ResponseMember
	if err := m.db.Model(&dao.Members{}).First(&responseMember, newMember.Id).Error; err != nil {
		return nil, err
	}

	return responseMember, nil
}

func NewMemberRepository(db *gorm.DB) MembersRepository {
	return &membersRepository{
		db: db,
	}
}
