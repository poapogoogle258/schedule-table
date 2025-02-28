package dao

import (
	"schedule_table/internal/constant"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	Id             uuid.UUID           `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId     uuid.UUID           `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ScheduleId     uuid.UUID           `gorm:"type:uuid;column:schedule_id" json:"schedule_id"`
	EmployeeId     *uuid.UUID          `gorm:"type:uuid;column:member_id" json:"member_id"`
	SecondMemberId *uuid.UUID          `gorm:"type:uuid;column:second_member_id" json:"second_member_id"`
	RecurrenceId   string              `gorm:"column:recurrence_id;size:56;unique" json:"recurrenceId"`
	Start          time.Time           `gorm:"column:start" json:"start"`
	End            time.Time           `gorm:"column:end" json:"end"`
	RestTime       time.Time           `gorm:"column:restTime" json:"restTime"`
	Priority       int8                `gorm:"column:priority" json:"priority"`
	Status         constant.TaskStatus `gorm:"column:status" json:"status"`
	Refreshed      time.Time           `gorm:"column:refreshed" json:"refreshed"`
	Person         *Employee           `gorm:"foreignKey:member_id" json:"person"`
	Description    Schedule            `gorm:"foreignKey:schedule_id" json:"description"`
	BaseModel
}

// update hook
func (task *Task) BeforeUpdate(tx *gorm.DB) error {
	if task.Status == constant.TaskCommitted {
		task.Status = constant.TaskCreated
	}

	return nil
}
func (task *Task) BeforeSave(tx *gorm.DB) error {
	if task.Status == constant.TaskCommitted {
		task.Status = constant.TaskCreated
	}

	return nil
}

// end:update hook
