package constant

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
