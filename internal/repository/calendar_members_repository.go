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
	CreateNewMember(newMember *dao.Members) (*dto.ResponseMember, error)
	EditMember(memberId string, insert *dao.Members) (*dto.ResponseMember, error)
	ExistMemberId(calendarId string, memberId string) bool
	DeleteMemberId(calendarId string, memberId string) error
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

func (m *membersRepository) CreateNewMember(newMember *dao.Members) (*dto.ResponseMember, error) {

	if err := m.db.Model(&dao.Members{}).Create(&newMember).Error; err != nil {
		return nil, err
	}

	var responseMember *dto.ResponseMember
	if err := m.db.Model(&dao.Members{}).First(&responseMember, newMember.Id).Error; err != nil {
		return nil, err
	}

	return responseMember, nil
}

func (m *membersRepository) EditMember(memberId string, insert *dao.Members) (*dto.ResponseMember, error) {

	if err := m.db.Model(&dao.Members{Id: util.Must(uuid.Parse(memberId))}).Updates(insert).Error; err != nil {
		return nil, err
	}

	var member *dto.ResponseMember
	if err := m.db.Model(&dao.Members{}).First(&member, `id = ?`, memberId).Error; err != nil {
		return nil, err
	}

	return member, nil

}

func (m *membersRepository) ExistMemberId(calendarId string, memberId string) bool {
	var countMember int64
	m.db.Model(&dao.Members{}).Where("id = ? AND calendar_id = ?", memberId, calendarId).Count(&countMember)

	return countMember > 0
}

func (m *membersRepository) DeleteMemberId(calendarId string, memberId string) error {
	return m.db.Delete(&dao.Members{}, "id = ? AND calendar_id = ?", memberId, calendarId).Error
}

func NewMemberRepository(db *gorm.DB) MembersRepository {
	return &membersRepository{
		db: db,
	}
}
