package dao

import (
	"schedule_table/internal/constant"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tasks struct {
	Id           uuid.UUID           `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	CalendarId   uuid.UUID           `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ScheduleId   uuid.UUID           `gorm:"type:uuid;column:schedule_id" json:"schedule_id"`
	MemberId     *uuid.UUID          `gorm:"type:uuid;column:member_id" json:"member_id"`
	RecurrenceId string              `gorm:"column:recurrence_id;size:49;unique" json:"recurrenceId"`
	Start        time.Time           `gorm:"column:start" json:"start"`
	End          time.Time           `gorm:"column:end" json:"end"`
	RestTime     time.Time           `gorm:"column:restTime" json:"restTime"`
	Priority     int8                `gorm:"column:priority" json:"priority"`
	Status       constant.TaskStatus `gorm:"column:status" json:"status"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	Person       *Members            `gorm:"foreignKey:member_id" json:"person"`
	Description  Schedules           `gorm:"foreignKey:schedule_id" json:"description"`
}

func (task *Tasks) BeforeCreate(tx *gorm.DB) (err error) {
	if task.Id == uuid.Nil {
		task.Id = uuid.New()
	}

	return
}
