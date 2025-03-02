package dto

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	Employee    EmployeeInfo `json:"employee"`
	TypeLeave   string       `json:"type"`
	Description string       `json:"description"`
	Date        time.Time    `json:"date"`
	Tzid        string       `json:"tzid"`
}

// dto -> dao
func (t *TaskRequest) EmployeeId() uuid.UUID {
	return uuid.MustParse(t.Employee.Id)
}

func (t *TaskRequest) Type() dao.LeaveType {
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

type TaskInfo struct {
	Id          string       `json:"id"`
	Employee    EmployeeInfo `json:"employee"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Date        time.Time    `json:"date"`
	Status      string       `json:"status"`
}
