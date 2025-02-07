package lib

import (
	"errors"
	"schedule_table/internal/constant"
	"schedule_table/internal/model/dao"
	"schedule_table/util"
	"time"

	"github.com/google/uuid"
)

type IWorker interface {
	GetId() uuid.UUID
	IsAvailable() bool
	AddTask(task *dao.Tasks) error
	// AddLeaven(day time.Time) error
	// AddReserved(task *dao.Tasks) error
}

var (
	ErrMemberNotAvailable             = errors.New("worker is not available")
	ErrMemberNotReady                 = errors.New("worker not ready")
	ErrMemberIsLeaven                 = errors.New("worker is leaven")
	ErrMemberReserved                 = errors.New("worker was reserved")
	ErrTaskIsReserved                 = errors.New("task is reserved")
	ErrNotMemberReserved              = errors.New("not use people who have been reserved")
	ErrTaskStatusNotCreatedOrReserved = errors.New("task status not created or reserved")
)

type Worker struct {
	Id            uuid.UUID
	Info          *dao.Members
	RestTime      *time.Time
	LeavesDays    []time.Time
	ReservedTasks []*dao.Tasks
}

func (worker *Worker) GetId() uuid.UUID {
	return worker.Id
}

func (worker *Worker) IsAvailable() bool {
	return worker.Info.Available
}

func (worker Worker) AddTask(task *dao.Tasks) error {

	if task.Status == constant.TaskStatus_Created {
		if !worker.IsAvailable() {
			return ErrMemberNotAvailable
		}

		if worker.RestTime != nil && task.Start.Before(*worker.RestTime) {
			return ErrMemberNotReady
		}

		for _, leaved := range worker.LeavesDays {
			if task.Start.Equal(leaved) || task.End.Equal(leaved) {
				return ErrMemberIsLeaven
			}
		}

		for _, reserve := range worker.ReservedTasks {
			if (task.Start.After(reserve.Start) && task.Start.Before(reserve.RestTime)) ||
				(task.End.After(reserve.Start) && task.End.Before(reserve.RestTime)) {
				return ErrTaskIsReserved
			}
		}
	} else if task.Status == constant.TaskStatus_Reserved {
		if *task.MemberId != worker.GetId() {
			return ErrNotMemberReserved
		}
	} else {
		return ErrTaskStatusNotCreatedOrReserved
	}

	worker.RestTime = &task.RestTime
	task.MemberId = &worker.Id
	task.Person = worker.Info

	return nil
}

func NewWorkerMember(member *dao.Members) IWorker {

	leavesDays := util.Map(member.Leaves, func(leave dao.Leaves) time.Time {
		return leave.Date
	})

	return &Worker{
		Id:         member.Id,
		Info:       member,
		LeavesDays: leavesDays,
	}
}
