package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tasks struct {
	Id          uuid.UUID  `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	CalendarId  uuid.UUID  `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ScheduleId  uuid.UUID  `gorm:"type:uuid;column:schedule_id" json:"schedule_id"`
	MemberId    uuid.UUID  `gorm:"type:uuid;column:member_id" json:"member_id"`
	Start       time.Time  `gorm:"column:start" json:"start"`
	End         time.Time  `gorm:"column:end" json:"end"`
	Priority    int8       `gorm:"column:priority" json:"priority"`
	Status      int8       `gorm:"column:status" json:"status"`
	Reserved    bool       `gorm:"column:reserved;default:false" json:"reserved"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Person      Members    `gorm:"foreignKey:member_id" json:"person"`
	Description *Schedules `gorm:"foreignKey:schedule_id" json:"description"`
}

func (task *Tasks) BeforeCreate(tx *gorm.DB) (err error) {
	task.Id = uuid.New()

	return
}
