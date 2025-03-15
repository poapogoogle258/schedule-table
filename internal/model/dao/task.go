package dao

import (
	"schedule_table/internal/constant"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id                 uuid.UUID           `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId         uuid.UUID           `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ScheduleId         uuid.UUID           `gorm:"type:uuid;column:schedule_id" json:"schedule_id"`
	EmployeeId         *uuid.UUID          `gorm:"type:uuid;column:employee_id" json:"employee_id"`
	OriginalEmployeeId *uuid.UUID          `gorm:"type:uuid;column:original_employee_id" json:"original_employee_id"`
	RecurrenceId       string              `gorm:"column:recurrence_id;size:56;unique" json:"recurrenceId"`
	Extra              bool                `gorm:"column:extra;default:false" json:"extra"`
	Start              time.Time           `gorm:"column:start" json:"start"`
	End                time.Time           `gorm:"column:end" json:"end"`
	RestTime           time.Time           `gorm:"column:restTime" json:"restTime"`
	Priority           int8                `gorm:"column:priority" json:"priority"`
	Status             constant.TaskStatus `gorm:"column:status" json:"status"`
	Refreshed          time.Time           `gorm:"column:refreshed" json:"refreshed"`
	Person             *Employee           `gorm:"foreignKey:employee_id" json:"person"`
	Description        Schedule            `gorm:"foreignKey:schedule_id" json:"description"`
	BaseModel
}

// end:update hook

func (task *Task) Committed() bool {
	return task.OriginalEmployeeId != nil
}
