package service

import (
	"errors"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"schedule_table/util"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type ScheduleService interface {
	GetAllSchedules(calendarId string) ([]*dto.ScheduleInfo, error)
	GetSchedule(scheduleId string) (*dto.ScheduleInfo, error)
	CreateSchedule(calendarId string, data *dto.ScheduleInfoRequest) (*dto.ScheduleInfo, error)
	UpdateSchedule(scheduleId string, data *dto.ScheduleInfoRequest) (*dto.ScheduleInfo, error)
	DeleteSchedule(scheduleId string) error
	IsExist(scheduleId string) bool
	IsExistOfCalendar(calendarId string, scheduleId string) bool
}

type scheduleService struct {
	transaction  repository.Transaction
	scheduleRepo repository.ScheduleRepository
}

func (s *scheduleService) GetAllSchedules(calendarId string) ([]*dto.ScheduleInfo, error) {
	schedules := util.Must(s.scheduleRepo.FindManyWithAggregateEmployee("calendar_id = ?", calendarId))

	result := []*dto.ScheduleInfo{}
	copier.Copy(&result, schedules)

	return result, nil
}

func (s *scheduleService) GetSchedule(scheduleId string) (*dto.ScheduleInfo, error) {
	schedule := util.Must(s.scheduleRepo.FindOneWithAggregateEmployee("id = ?", scheduleId))

	result := &dto.ScheduleInfo{}
	copier.Copy(&result, schedule)

	return result, nil
}

var ErrMasterScheduleNotFound = errors.New("master schedule not found")

func (s *scheduleService) CreateSchedule(calendarId string, data *dto.ScheduleInfoRequest) (*dto.ScheduleInfo, error) {

	insert := &dao.Schedule{}
	copier.Copy(&insert, data)

	insert.Id = uuid.New()
	insert.CalendarId = uuid.MustParse(calendarId)
	if len(data.MasterScheduleId) == 0 {
		insert.MasterScheduleId = insert.Id
	} else {
		if !s.scheduleRepo.IsExist("calendar_id = ?", insert.MasterScheduleId) {
			return nil, ErrMasterScheduleNotFound
		}
	}

	tx := s.transaction.Begin()
	defer tx.Rollback()

	scheduleRepo := s.scheduleRepo.InjectionTx(tx)

	if err := scheduleRepo.Create(insert); err != nil {
		return nil, err
	}

	employeeIds := make([]string, len(data.Employees))
	for i := range employeeIds {
		employeeIds[i] = data.Employees[i].Id
	}
	if err := scheduleRepo.UpdateEmployeeQueue(insert.Id.String(), employeeIds); err != nil {
		return nil, err
	}

	tx.Commit()
	return s.GetSchedule(insert.Id.String())

}

func (s *scheduleService) UpdateSchedule(scheduleId string, data *dto.ScheduleInfoRequest) (*dto.ScheduleInfo, error) {

	schedule := util.Must(s.scheduleRepo.FindOne("id = ?", scheduleId))
	copier.Copy(&schedule, data)

	if err := s.scheduleRepo.Save(schedule); err != nil {
		return nil, err
	}

	employeeIds := make([]string, len(data.Employees))
	for i := range employeeIds {
		employeeIds[i] = data.Employees[i].Id
	}
	if err := s.scheduleRepo.UpdateEmployeeQueue(scheduleId, employeeIds); err != nil {
		return nil, err
	}

	return s.GetSchedule(scheduleId)

}

func (s *scheduleService) DeleteSchedule(scheduleId string) error {

	tx := s.transaction.Begin()
	defer tx.Rollback()

	scheduleRepo := s.scheduleRepo.InjectionTx(tx)

	if err := scheduleRepo.ClearEmployeeQueue(scheduleId); err != nil {
		return err
	}

	if err := scheduleRepo.Delete("id = ?", scheduleId); err != nil {
		return err
	}

	tx.Commit()
	return nil

}

func (s *scheduleService) IsExist(scheduleId string) bool {
	return s.scheduleRepo.IsExist("id = ?", scheduleId)
}

func (s *scheduleService) IsExistOfCalendar(calendarId string, scheduleId string) bool {
	return s.scheduleRepo.IsExist("calendar_id = ? AND id = ?", calendarId, scheduleId)
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, transaction repository.Transaction) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		transaction:  transaction,
	}
}
