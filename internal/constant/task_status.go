package constant

import "errors"

type TaskStatus int8

const (
	TaskStatus_Created TaskStatus = iota
	TaskStatus_Submitted
	TaskStatus_Reserved
	TaskStatus_Canceled
)

func (status TaskStatus) ToString() string {
	switch status {
	case TaskStatus_Created:
		return "Created"
	case TaskStatus_Submitted:
		return "Submitted"
	case TaskStatus_Reserved:
		return "Reserved"
	case TaskStatus_Canceled:
		return "Canceled"
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
