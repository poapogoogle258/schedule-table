package repository

import (
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Repository[*dao.Task]
	InjectionTx(db *gorm.DB) TaskRepository
	FindManyWithJoin(conds ...any) ([]*dao.Task, error)
	FindOneWithJoin(conds ...any) (*dao.Task, error)
	Upsert(tasks []*dao.Task) error
}

type taskRepository struct {
	Repository[*dao.Task]
	db *gorm.DB
}

var (
	fieldEmployeeInfoShort = []string{"id", "name", "nickname", "position", "image_url"}
	fieldScheduleInfoShort = []string{"id", "name", "image_url", "color"}
)

func (repo *taskRepository) Upsert(tasks []*dao.Task) error {
	updatedAt := time.Now()

	newTasks := make([]*dao.Task, 0)
	oldTasks := make([]*dao.Task, 0)

	for i := range tasks {
		tasks[i].Refreshed = updatedAt
		if tasks[i].Status == constant.TaskGenerated {
			tasks[i].Status = constant.TaskCreated
			newTasks = append(newTasks, tasks[i])
		} else {
			oldTasks = append(oldTasks, tasks[i])
		}
	}

	// upsert data
	return repo.db.Transaction(func(tx *gorm.DB) error {

		if len(newTasks) > 0 {
			if err := tx.Create(newTasks).Error; err != nil {
				return err
			}
		}

		if len(oldTasks) > 0 {
			if err := tx.Save(oldTasks).Error; err != nil {
				return err
			}
		}

		return nil
	})

}

func (repo *taskRepository) FindOneWithJoin(conds ...any) (*dao.Task, error) {
	var task *dao.Task

	result := repo.db.
		Joins("Person", repo.db.Select(fieldEmployeeInfoShort)).
		Joins("Description", repo.db.Select(fieldScheduleInfoShort)).
		First(&task, conds...)

	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repo *taskRepository) FindManyWithJoin(conds ...any) ([]*dao.Task, error) {
	var tasks []*dao.Task

	result := repo.db.
		Joins("Person", repo.db.Select(fieldEmployeeInfoShort)).
		Joins("Description", repo.db.Select(fieldScheduleInfoShort)).
		Find(&tasks, conds...)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (repo *taskRepository) InjectionTx(db *gorm.DB) TaskRepository {
	return NewTaskRepository(db)
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		Repository: NewRepository[*dao.Task](db),
		db:         db,
	}
}
