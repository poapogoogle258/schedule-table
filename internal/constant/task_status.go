package constant

type TaskStatus int8

const (
	TaskGenerated TaskStatus = iota // is new task
	TaskCreated
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

// var ErrTaskStatusNotExist = errors.New("task status not exist")
