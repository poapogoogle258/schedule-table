package repository

import (
	"errors"
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	Find(conds ...interface{}) (*[]dao.Tasks, error)
	First(conds ...interface{}) (*dao.Tasks, error)
	FindOrderLimit(order string, limit int, conds ...interface{}) (*[]dao.Tasks, error)
	UpdatesAndFind(taskId string, value interface{}) (*dao.Tasks, error)
	CreateTasks(insert []*dao.Tasks) error
	DeleteOne(taskId uuid.UUID) error
	FindWithAssociation(conds ...interface{}) (*[]dao.Tasks, error)
	IsExists(taskId string) bool
}

type TaskRepository struct {
	db *gorm.DB
}

var (
	ErrTaskNotExists = errors.New("task not exits")
)

func (taskRepo *TaskRepository) First(conds ...interface{}) (*dao.Tasks, error) {
	task := &dao.Tasks{}
	if err := taskRepo.db.First(&task, conds...).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (taskRepo *TaskRepository) IsExists(taskId string) bool {
	var count int64
	if err := taskRepo.db.Model(&dao.Tasks{}).Where("id = ?", taskId).Limit(1).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (taskRepo *TaskRepository) FindWithAssociation(conds ...interface{}) (*[]dao.Tasks, error) {
	tasks := &[]dao.Tasks{}

	if err := taskRepo.db.Preload("Description").Preload("Person").Find(&tasks, conds...).Error; err != nil {
		return nil, err
	}

	return tasks, nil

}

func (taskRepo *TaskRepository) CreateTasks(insert []*dao.Tasks) error {

	if err := taskRepo.db.Model(&dao.Tasks{}).Create(insert).Error; err != nil {
		return err
	}

	return nil
}

func (taskRepo *TaskRepository) Find(conds ...interface{}) (*[]dao.Tasks, error) {
	var tasks *[]dao.Tasks
	if err := taskRepo.db.Find(&tasks, conds...).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) FindOrderLimit(order string, limit int, conds ...interface{}) (*[]dao.Tasks, error) {
	var tasks *[]dao.Tasks
	if err := taskRepo.db.Limit(limit).Order(order).Find(&tasks, conds...).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TaskRepository) UpdatesAndFind(taskId string, value interface{}) (*dao.Tasks, error) {
	task := &dao.Tasks{}

	if err := repo.db.Model(&task).Clauses(clause.Returning{}).Where("id = ?", taskId).Updates(value).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (taskRepo *TaskRepository) Delete(order string, limit int, conds ...interface{}) (*[]dao.Tasks, error) {
	var tasks *[]dao.Tasks
	if err := taskRepo.db.Limit(limit).Order(order).Find(&tasks, conds...).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) DeleteOne(taskId uuid.UUID) error {
	return taskRepo.db.Delete(&dao.Tasks{}, "id = ?", taskId).Error
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}
