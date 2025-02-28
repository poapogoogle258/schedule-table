package lib

import (
	"errors"
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type WorkType int8

const (
	Works WorkType = iota
	Rest
	Reserved
	Leave
	Submitted
)

var (
	ErrTaskDataInvalid  = errors.New("task invalid")
	ErrLeaveDataInvalid = errors.New("leave invalid")
	ErrWortTypeInvalid  = errors.New("work type invalid")
)

type Work struct {
	Id          uuid.UUID
	ReferenceId *uuid.UUID
	Type        WorkType
	Priority    int8
	Start       time.Time
	End         time.Time
}

func NewWork(id uuid.UUID, referId *uuid.UUID, start time.Time, end time.Time, priority int8, workType WorkType) *Work {
	return &Work{
		Id:          id,
		ReferenceId: referId,
		Start:       start,
		End:         end,
		Priority:    priority,
		Type:        workType,
	}
}

func FactoryWork(data interface{}, workType WorkType) (*Work, error) {

	switch workType {
	case Works:
		task, ok := data.(*dao.Task)
		if !ok {
			return nil, ErrTaskDataInvalid
		}
		return NewWork(task.Id, nil, task.Start, task.End, task.Priority, Works), nil
	case Rest:
		task, ok := data.(*dao.Task)
		if !ok {
			return nil, ErrTaskDataInvalid
		}
		return NewWork(uuid.New(), &task.Id, task.End, task.End.Add(time.Minute*time.Duration(task.Description.BreakTime)), task.Priority, Rest), nil
	case Reserved:
		task, ok := data.(*dao.Task)
		if !ok {
			return nil, ErrTaskDataInvalid
		}
		return NewWork(task.Id, nil, task.Start, task.End, task.Priority, Reserved), nil
	case Submitted:
		task, ok := data.(*dao.Task)
		if !ok {
			return nil, ErrTaskDataInvalid
		}
		return NewWork(task.Id, nil, task.Start, task.End, task.Priority, Submitted), nil
	case Leave:
		leave, ok := data.(*dao.Leave)
		if !ok {
			return nil, ErrLeaveDataInvalid
		}

		location, err := time.LoadLocation(leave.Tzid)
		if err != nil {
			panic(err)
		}
		leaveDate := leave.Date.In(location)
		start := time.Date(leaveDate.Year(), leaveDate.Month(), leaveDate.Day(), 0, 0, 0, 0, location)
		end := time.Date(leaveDate.Year(), leaveDate.Month(), leaveDate.Day(), 23, 59, 59, 0, location)

		return NewWork(leave.Id, nil, start, end, 0, Leave), nil
	default:
		return nil, ErrWortTypeInvalid
	}

}
