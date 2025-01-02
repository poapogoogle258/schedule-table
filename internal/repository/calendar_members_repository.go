package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type MembersRepository interface {
	GetMemberId(memberId string) (*dao.Members, error)
	GetMembers(calendarId string) (*[]dao.Members, error)
	CreateNewMember(newMember *dao.Members) (*dao.Members, error)
	EditMember(memberId string, insert *dao.Members) (*dao.Members, error)
	ExistMemberId(calendarId string, memberId string) bool
	DeleteMemberId(calendarId string, memberId string) error
}

type membersRepository struct {
	db *gorm.DB
}

func (m *membersRepository) GetMemberId(memberId string) (*dao.Members, error) {
	var member *dao.Members

	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}

	if err := m.db.Model(&dao.Members{}).Select(selectedField).First(&member, "id = ?", memberId).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (m *membersRepository) GetMembers(calendarId string) (*[]dao.Members, error) {
	var members *[]dao.Members
	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}

	if err := m.db.Model(&dao.Members{}).Select(selectedField).Find(&members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (m *membersRepository) CreateNewMember(newMember *dao.Members) (*dao.Members, error) {

	if err := m.db.Model(&dao.Members{}).Create(&newMember).Error; err != nil {
		return nil, err
	}

	return newMember, nil
}

func (m *membersRepository) EditMember(memberId string, insert *dao.Members) (*dao.Members, error) {
	var member *dao.Members
	if err := m.db.Model(&dao.Members{}).First(&member, "id = ?", memberId).Error; err != nil {
		return nil, err
	}

	if err := m.db.Model(&member).Updates(insert).Error; err != nil {
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
