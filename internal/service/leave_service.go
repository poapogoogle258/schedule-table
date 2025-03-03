package service

import (
	"errors"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type LeaveService interface {
	IsExist(taskId string) bool
	IsExistOfCalendar(calendarId string, taskId string) bool
	FindAllLeavesOfCalendar(calendarId string, start time.Time, end time.Time) ([]*dto.LeaveInfo, error)
	FindOndLeave(id string) (*dto.LeaveInfo, error)
	CreateLeave(calendarId string, request *dto.LeaveRequest) (*dto.LeaveInfo, error)
	ChangeStatusLeave(leaveId string, newStatus string) error
}

type leaveService struct {
	leaveRepo repository.LeaveRepository
}

func (s *leaveService) IsExistOfCalendar(calendarId string, taskId string) bool {
	return s.leaveRepo.IsExist("id = ? AND calendar_id = ?", calendarId, taskId)
}

func (s *leaveService) FindAllLeavesOfCalendar(calendarId string, start time.Time, end time.Time) ([]*dto.LeaveInfo, error) {
	leaves, err := s.leaveRepo.FindManyWithAggregate([]string{"Employee"}, "calendar_id = ? AND start_time BETWEEN ? AND ?", calendarId, start, end)
	if err != nil {
		return nil, err
	}

	results := []*dto.LeaveInfo{}
	copier.Copy(&results, leaves)

	return results, nil
}

func (s *leaveService) FindOndLeave(id string) (*dto.LeaveInfo, error) {
	leave, err := s.leaveRepo.FindOneWithAggregate([]string{"Employee"}, "id = ?", id)
	if err != nil {
		return nil, err
	}

	result := &dto.LeaveInfo{}
	copier.Copy(result, leave)

	return result, nil
}

func (s *leaveService) IsExist(taskId string) bool {
	return s.leaveRepo.IsExist("id = ?", taskId)
}

func (s *leaveService) CreateLeave(calendarId string, request *dto.LeaveRequest) (*dto.LeaveInfo, error) {
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

var (
	ErrLeaveSameStatus    = errors.New("status is same")
	ErrLeaveStatusNotFlow = errors.New("reject status not flow state")
	ErrInvalidStatus      = errors.New("invalid leave status")
)

func (s *leaveService) ChangeStatusLeave(leaveId string, newStatus string) error {

	mapLeaveStatus := map[string]dao.LeaveStatus{
		"pending": dao.Pending,
		"accept":  dao.Accept,
		"reject":  dao.Reject,
		"cancel":  dao.Cancel,
	}

	if _, ok := mapLeaveStatus[newStatus]; !ok {
		return ErrInvalidStatus
	}

	leave, err := s.leaveRepo.FindOne("id = ?", leaveId)
	if err != nil {
		return err
	}

	// check state flow
	switch leave.Status {
	case dao.Pending:
		// padding -> Accept, Reject, Cancel
		if newStatus == "accept" || newStatus == "reject" || newStatus == "cancel" {
			break
		} else {
			return ErrLeaveStatusNotFlow
		}
	case dao.Accept:
		// Accept -> Cancel
		if newStatus == "cancel" {
			break
		} else {
			return ErrLeaveStatusNotFlow
		}
	case dao.Reject:
		return ErrLeaveStatusNotFlow
	case dao.Cancel:
		return ErrLeaveStatusNotFlow
	default:
		return ErrInvalidStatus
	}

	leave.Status = mapLeaveStatus[newStatus]
	updateAt := time.Now()
	switch leave.Status {
	case dao.Accept:
		leave.AcceptAt = &updateAt
	case dao.Reject:
		leave.RejectAt = &updateAt
	case dao.Cancel:
		leave.CancelAt = &updateAt
	}

	if err := s.leaveRepo.Update(leave); err != nil {
		return err
	}

	return nil

}

func NewLeaveService(leaveRepo repository.LeaveRepository) LeaveService {
	return &leaveService{
		leaveRepo: leaveRepo,
	}
}
