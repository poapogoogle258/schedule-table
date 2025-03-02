package service

import (
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type LeaveService interface {
	IsExist(taskId string) bool
	FindAllLeavesOfCalendar(calendarId string, start time.Time, end time.Time) ([]*dto.TaskInfo, error)
	FindOndLeave(id string) (*dto.TaskInfo, error)
	CreateLeave(calendarId string, request *dto.TaskRequest) (*dto.TaskInfo, error)
}

type leaveService struct {
	leaveRepo repository.LeaveRepository
}

func (s *leaveService) FindAllLeavesOfCalendar(calendarId string, start time.Time, end time.Time) ([]*dto.TaskInfo, error) {
	leaves, err := s.leaveRepo.FindManyWithAggregate([]string{"Employee"}, "calendar_id = ? AND start_time BETWEEN ? AND ?", calendarId, start, end)
	if err != nil {
		return nil, err
	}

	results := []*dto.TaskInfo{}
	copier.Copy(&results, leaves)

	return results, nil
}

func (s *leaveService) FindOndLeave(id string) (*dto.TaskInfo, error) {
	leave, err := s.leaveRepo.FindOneWithAggregate([]string{"Employee"}, "id = ?", id)
	if err != nil {
		return nil, err
	}

	result := &dto.TaskInfo{}
	copier.Copy(result, leave)

	return result, nil
}

func (s *leaveService) IsExist(taskId string) bool {
	return s.leaveRepo.IsExist("id = ?", taskId)
}

func (s *leaveService) CreateLeave(calendarId string, request *dto.TaskRequest) (*dto.TaskInfo, error) {
	insert := &dao.Leave{}
	copier.Copy(&insert, request)

	insert.Id = uuid.New()
	insert.CalendarId = uuid.MustParse(calendarId)
	insert.Status = dao.Pending

	if err := s.leaveRepo.Create(insert); err != nil {
		return nil, err
	}

	return s.FindOndLeave(insert.Id.String())

}

func NewLeaveServiceImpl(leaveRepo repository.LeaveRepository) LeaveService {
	return &leaveService{
		leaveRepo: leaveRepo,
	}
}
