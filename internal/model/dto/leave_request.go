package dto

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type LeaveRequest struct {
	Employee    EmployeeInfo `json:"employee" binding:"required"`
	TypeLeave   string       `json:"type" binding:"required"`
	Description string       `json:"description" binding:"required"`
	DateOnly    string       `json:"date" binding:"required" time_format:"2006-01-02"`
	Tzid        string       `json:"tzid" binding:"required"`
}

// dto -> dao
func (t *LeaveRequest) EmployeeId() uuid.UUID {
	return uuid.MustParse(t.Employee.Id)
}

func (t *LeaveRequest) Date() time.Time {
	date, err := time.Parse("2006-01-02", t.DateOnly)
	if err != nil {
		panic(err)
	}
	return date
}

func (t *LeaveRequest) Type() dao.LeaveType {
	switch t.TypeLeave {

	case "Sick":
		return dao.SickLeaveType
	case "Personal":
		return dao.PersonalLeaveType
	case "Annual":
		return dao.AnnualLeaveType
	default:
		panic("Unknown leave type")
	}
}

type LeaveInfo struct {
	Id           string       `json:"id"`
	Employee     EmployeeInfo `json:"employee"`
	Type         string       `json:"type"`
	Description  string       `json:"description"`
	Date         time.Time    `json:"date"`
	StatusString string       `json:"status"`
}

func (s *LeaveInfo) Status(v dao.LeaveStatus) {
	s.StatusString = v.String()
}
