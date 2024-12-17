package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Responsible struct {
	gorm.Model
	MemberId     uuid.UUID  `gorm:"type:uuid;column:member_id;primaryKey" json:"member_id"`
	ScheduleId   uuid.UUID  `gorm:"type:uuid;column:schedule_id;primaryKey" json:"schedule_id"`
	Queue        int8       `gorm:"column:queue" json:"queue"`
	Limit        int8       `gorm:"column:limit" json:"limit"`
	LastTimeTask *time.Time `gorm:"column:lastTimeTask" json:"lastTimeTask"`
	Person       Members    `gorm:"foreignKey:member_id" json:"person"`
}
