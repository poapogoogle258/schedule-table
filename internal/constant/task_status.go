package constant

import "errors"

type TaskStatus int8

const (
	TaskGenerated TaskStatus = iota
	TaskCreated
	TaskCommitted
	TaskReserved
	TaskOvertime
	TaskCanceled
)

func (status TaskStatus) ToString() string {
	switch status {
	case TaskGenerated:
		return "Generated"
	case TaskCreated:
		return "Created"
	case TaskCommitted:
		return "Submitted"
	case TaskReserved:
		return "Reserved"
	case TaskOvertime:
		return "Overtime"
	case TaskCanceled:
		return "Cancel"
	default:
		panic("Not Definition TaskStatus")
	}
}

var ErrTaskStatusNotExist = errors.New("task status not exist")

func ParseTaskStatus(status int8) (TaskStatus, error) {
	if status == 0 || status == 1 || status == 2 || status == 3 {
		return TaskStatus(status), nil
	}

	return 0, ErrTaskStatusNotExist
}
