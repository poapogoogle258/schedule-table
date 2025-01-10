package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	Find(conds ...interface{}) (*[]dao.Tasks, error)
	UpdatesAndFind(taskId string, value interface{}) (*dao.Tasks, error)
}

type TaskRepository struct {
	db *gorm.DB
}

func (taskRepo *TaskRepository) Find(conds ...interface{}) (*[]dao.Tasks, error) {
	var tasks *[]dao.Tasks
	if err := taskRepo.db.Find(&tasks, conds...).Error; err != nil {
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

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}
