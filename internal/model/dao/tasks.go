package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	Id          uuid.UUID `gorm:"type:uuid; column:id; primary_key; uniqueIndex" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ScheduleId  uuid.UUID `gorm:"type:uuid;column:schedule_id" json:"schedule_id"`
	MemberId    uuid.UUID `gorm:"type:uuid;column:member_id" json:"member_id"`
	Start       time.Time `gorm:"column:start" json:"start"`
	End         time.Time `gorm:"column:end" json:"end"`
	Status      int8      `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Person      Members   `gorm:"foreignKey:member_id" json:"person"`
	Description Schedules `gorm:"foreignKey:schedule_id" json:"description"`
}
