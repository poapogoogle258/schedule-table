package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Calendars struct {
	Id                           uuid.UUID    `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	Name                         string       `gorm:"column:name" json:"name"`
	ImageURL                     string       `gorm:"column:image_url;default:default-member-profile.jpeg" json:"imageURL"`
	Description                  string       `gorm:"column:description" json:"description"`
	UserId                       uuid.UUID    `gorm:"type:uuid;column:user_id" json:"user_id"`
	LastTimeUpdatedRecurrence    time.Time    `gorm:"column:updated_recurrence" json:"updated_recurrence"`
	LastTimeUpdatedGenerateTasks *time.Time   `gorm:"column:updated_generate_tasks" json:"updated_generate_tasks"`
	Members                      *[]Members   `gorm:"foreignKey:calendar_id" json:"members"`
	Leaves                       *[]Leaves    `gorm:"foreignKey:calendar_id" json:"leaves"`
	Schedules                    *[]Schedules `gorm:"foreignKey:calendar_id" json:"schedules"`
	Tasks                        *[]Tasks     `gorm:"foreignKey:calendar_id" json:"tasks"`
	BaseModel
}

func (cal *Calendars) BeforeCreate(tx *gorm.DB) (err error) {
	if cal.Id == uuid.Nil {
		cal.Id = uuid.New()
	}

	return
}
