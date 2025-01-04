package repository

import (
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type MembersRepository interface {
	FindOne(memberId string) (*dao.Members, error)
	FindByCalendarId(calendarId string) (*[]dao.Members, error)
	Create(newMember *dao.Members) (*dao.Members, error)
	UpdateOne(memberId string, insert *dao.Members) (*dao.Members, error)
	IsExits(memberId string) bool
	DeleteOne(memberId string) error
}

type membersRepository struct {
	db *gorm.DB
}

func selectColumnMember(db *gorm.DB) *gorm.DB {
	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}
	return db.Select(selectedField)
}

func (memRepo *membersRepository) FindOne(memberId string) (*dao.Members, error) {
	var member *dao.Members

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).First(&member, "id = ?", memberId).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (memRepo *membersRepository) FindByCalendarId(calendarId string) (*[]dao.Members, error) {
	var members *[]dao.Members

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).Find(&members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (memRepo *membersRepository) Create(newMember *dao.Members) (*dao.Members, error) {

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).Create(&newMember).Error; err != nil {
		return nil, err
	}

	return newMember, nil
}

func (memRepo *membersRepository) UpdateOne(memberId string, insert *dao.Members) (*dao.Members, error) {

	member := &dao.Members{}
	copier.Copy(&member, &insert)

	member.Id = uuid.MustParse(memberId)

	if err := memRepo.db.Scopes(selectColumnMember).Updates(&member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (memRepo *membersRepository) IsExits(memberId string) bool {
	var countMember int64
	if err := memRepo.db.Model(&dao.Members{}).Where("id = ?", memberId).Count(&countMember).Error; err != nil {
		panic(err)
	}

	return countMember > 0
}

func (memRepo *membersRepository) DeleteOne(memberId string) error {
	return memRepo.db.Delete(&dao.Members{}, "id = ?", memberId).Error
}

func NewMemberRepository(db *gorm.DB) MembersRepository {
	return &membersRepository{
		db: db,
	}
}
