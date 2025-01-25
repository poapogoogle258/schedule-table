package repository

import (
	"errors"
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrMemberNotFount = errors.New("not fount member")
)

type MembersRepository interface {
	FindOne(conds ...interface{}) (*dao.Members, error)
	Find(conds ...interface{}) (*[]dao.Members, error)
	FindWithOffsetAndLimit(offset int, limit int, conds ...interface{}) (*[]dao.Members, error)
	Count(calendarId string) int64
	Create(newMember *dao.Members) error
	UpdatesAndFindOne(memberId string, calendarId string, insert *dao.Members) (*dao.Members, error)
	DeleteOne(memberId string, calendarId string) error
	CheckExist(memberId string) error
}

type membersRepository struct {
	db *gorm.DB
}

func selectColumnMember(db *gorm.DB) *gorm.DB {
	selectedField := []string{"Id", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone"}
	return db.Select(selectedField)
}

func (memRepo *membersRepository) FindOne(conds ...interface{}) (*dao.Members, error) {
	var member *dao.Members

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).First(&member, conds...).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (memRepo *membersRepository) Find(conds ...interface{}) (*[]dao.Members, error) {
	var members *[]dao.Members

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).Find(&members, conds...).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (memRepo *membersRepository) FindWithOffsetAndLimit(offset int, limit int, conds ...interface{}) (*[]dao.Members, error) {
	var members *[]dao.Members

	if err := memRepo.db.Model(&dao.Members{}).Scopes(selectColumnMember).Offset(offset).Limit(limit).Find(&members, conds...).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (memRepo *membersRepository) Create(newMember *dao.Members) error {

	if err := memRepo.db.Model(&dao.Members{}).Create(&newMember).Error; err != nil {
		return err
	}

	return nil
}

func (memRepo *membersRepository) UpdatesAndFindOne(memberId string, calendarId string, insert *dao.Members) (*dao.Members, error) {

	member := &dao.Members{}
	if err := memRepo.db.First(&member, map[string]interface{}{
		"id":          memberId,
		"calendar_id": calendarId,
	}).Error; err != nil {
		return nil, err
	}

	if err := memRepo.db.Model(&member).Clauses(clause.Returning{}).Updates(insert).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (memRepo *membersRepository) CheckExist(memberId string) error {
	var countMember int64
	if err := memRepo.db.Model(&dao.Members{}).Where("id = ?", memberId).Count(&countMember).Error; err != nil {
		panic(err)
	}

	if countMember == 0 {
		return ErrMemberNotFount
	} else {
		return nil
	}

}

func (memRepo *membersRepository) DeleteOne(memberId string, calendarId string) error {
	return memRepo.db.Delete(&dao.Members{}, "id = ? AND calendar_id = ?", memberId, calendarId).Error
}

func (memRepo *membersRepository) Count(calendarId string) int64 {
	var count int64
	if err := memRepo.db.Model(&dao.Members{}).Where("calendar_id = ?", calendarId).Count(&count).Error; err != nil {
		panic(err)
	}

	return count
}

func NewMemberRepository(db *gorm.DB) MembersRepository {
	return &membersRepository{
		db: db,
	}
}
