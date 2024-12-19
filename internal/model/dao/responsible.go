package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Responsible struct {
	gorm.Model
	MemberId   uuid.UUID `gorm:"type:uuid;column:member_id;primaryKey" json:"member_id"`
	ScheduleId uuid.UUID `gorm:"type:uuid;column:schedule_id;primaryKey" json:"schedule_id"`
	Queue      int8      `gorm:"column:queue" json:"queue"`
	Person     Members   `gorm:"foreignKey:member_id" json:"person"`
}
