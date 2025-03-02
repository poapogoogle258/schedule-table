package repository

import (
	"schedule_table/internal/model/dao"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Repository[*dao.Schedule]
	UpdateEmployeeQueue(scheduleId string, employeeIds []string) error
	ClearEmployeeQueue(scheduleId string) error
	FindOneWithAggregateEmployee(conds ...any) (*dao.Schedule, error)
	FindManyWithAggregateEmployee(conds ...any) ([]*dao.Schedule, error)
	InjectionTx(tx *gorm.DB) ScheduleRepository
}

type scheduleRepository struct {
	Repository[*dao.Schedule]
	db *gorm.DB
}

func (repo *scheduleRepository) InjectionTx(tx *gorm.DB) ScheduleRepository {
	return &scheduleRepository{
		Repository: NewRepository[*dao.Schedule](tx),
		db:         tx,
	}
}

func (repo *scheduleRepository) ClearEmployeeQueue(scheduleId string) error {
	return repo.db.Where("schedule_id = ?", scheduleId).Delete(&dao.EmployeeQueue{}).Error
}

func (repo *scheduleRepository) UpdateEmployeeQueue(scheduleId string, employeeIds []string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// delete old queue
		if err := tx.Delete(&dao.EmployeeQueue{}, "schedule_id = ?", scheduleId).Error; err != nil {
			return err
		}

		// create
		insert := make([]*dao.EmployeeQueue, len(employeeIds))
		for i, employeeId := range employeeIds {
			insert[i] = &dao.EmployeeQueue{
				ScheduleId: uuid.MustParse(scheduleId),
				EmployeeId: uuid.MustParse(employeeId),
				Queue:      int8(i + 1),
			}
		}

		if err := tx.Create(insert).Error; err != nil {
			return err
		}

		return nil
	})

}

func (repo *scheduleRepository) FindOneWithAggregateEmployee(conds ...any) (*dao.Schedule, error) {

	var schedule *dao.Schedule

	if err := repo.db.First(&schedule, conds...).Error; err != nil {
		return nil, err
	}

	var employees []*dao.EmployeeQueue

	if err := repo.db.Joins("Person").Order("queue ASC").Find(&employees, "schedule_id = ?", schedule.MasterScheduleId).Error; err != nil {
		return nil, err
	}

	schedule.EmployeeQueue = employees

	return schedule, nil

}

func (repo *scheduleRepository) FindManyWithAggregateEmployee(conds ...any) ([]*dao.Schedule, error) {

	var schedules []*dao.Schedule

	if err := repo.db.Find(&schedules, conds...).Error; err != nil {
		return nil, err
	}

	var employees []*dao.EmployeeQueue
	mapScheduleMasterIds := make(map[string]bool)
	for _, schedule := range schedules {
		mapScheduleMasterIds[schedule.MasterScheduleId.String()] = true
	}

	if err := repo.db.Joins("Person").Order("queue ASC").Find(&employees, "schedule_id in (?)", maps.Keys(mapScheduleMasterIds)).Error; err != nil {
		return nil, err
	}

	groupByMasterId := make(map[uuid.UUID][]*dao.EmployeeQueue)
	for i := range employees {
		groupByMasterId[employees[i].ScheduleId] = append(groupByMasterId[employees[i].ScheduleId], employees[i])
	}

	for i := range schedules {
		schedules[i].EmployeeQueue = groupByMasterId[schedules[i].MasterScheduleId]
	}

	return schedules, nil

}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{
		db:         db,
		Repository: NewRepository[*dao.Schedule](db),
	}
}
