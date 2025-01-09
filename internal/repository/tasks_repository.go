package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type ITaskRepository interface {
	Find(conds ...interface{}) (*[]dao.Tasks, error)
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

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}
