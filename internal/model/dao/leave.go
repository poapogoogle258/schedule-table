package dao

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType string

const (
	SickLeaveType     LeaveType = "Sick"
	PersonalLeaveType LeaveType = "Personal"
	AnnualLeaveType   LeaveType = "Annual"
)

type LeaveStatus uint8

const (
	Pending LeaveStatus = iota
	Accept
	Reject
	Cancel
)

func (status LeaveStatus) String() string {
	switch status {
	case Pending:
		return "Pending"
	case Accept:
		return "Accept"
	case Reject:
		return "Reject"
	case Cancel:
		return "Cancel"
	default:
		return "Unknown"
	}
}

type Leave struct {
	Id          uuid.UUID   `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId  uuid.UUID   `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	EmployeeId  uuid.UUID   `gorm:"type:uuid;column:employee_id" json:"employee_id"`
	Type        LeaveType   `gorm:"column:type;size:10" json:"type"`
	Description string      `gorm:"column:description" json:"description"`
	Date        time.Time   `gorm:"column:date" json:"date"`
	Tzid        string      `gorm:"column:tzid;default:Asia/Bangkok" json:"tzid"`
	Status      LeaveStatus `gorm:"column:status;default:0" json:"status"`
	Employee    Employee    `gorm:"foreignKey:employee_id"`
	AcceptAt    *time.Time  `gorm:"column:accept_at" json:"accept_at"`
	RejectAt    *time.Time  `gorm:"column:reject_at" json:"reject_at"`
	CancelAt    *time.Time  `gorm:"column:cancel_at" json:"cancel_at"`
	BaseModel
}
