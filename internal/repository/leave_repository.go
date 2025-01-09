package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type LeaveRepository interface {
	Find(conds ...interface{}) (*[]dao.Leaves, error)
	FindOne(conds ...interface{}) (*dao.Leaves, error)
	Create(insert *dao.Leaves) error
	Delete(leaveId string) error
	Exits(leaveId string) bool
}

type leaveRepository struct {
	db *gorm.DB
}

func (leaveRepo *leaveRepository) Find(conds ...interface{}) (*[]dao.Leaves, error) {
	var leaves *[]dao.Leaves
	if err := leaveRepo.db.Joins("Member").Find(&leaves, conds...).Error; err != nil {
		return nil, err
	}

	return leaves, nil
}

func (leaveRepo *leaveRepository) FindOne(conds ...interface{}) (*dao.Leaves, error) {
	var leave *dao.Leaves
	if err := leaveRepo.db.Joins("Member").First(&leave, conds...).Error; err != nil {
		return nil, err
	}

	return leave, nil
}

func (leaveRepo *leaveRepository) Create(insert *dao.Leaves) error {
	if err := leaveRepo.db.Create(insert).Error; err != nil {
		return err
	}

	return nil
}

func (leaveRepo *leaveRepository) Exits(leaveId string) bool {
	var count int64
	if err := leaveRepo.db.Model(&dao.Leaves{}).Where("id = ?", leaveId).Count(&count).Error; err != nil {
		panic(err)
	}

	return count > 0
}

func (leaveRepo *leaveRepository) Delete(leaveId string) error {
	if err := leaveRepo.db.Delete(&dao.Leaves{}, "id = ?", leaveId).Error; err != nil {
		return err
	}

	return nil
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{
		db: db,
	}
}
