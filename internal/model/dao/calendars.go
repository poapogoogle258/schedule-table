package dao

import (
	"schedule_table/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Calendars struct {
	Id                      uuid.UUID    `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	Name                    string       `gorm:"column:name" json:"name"`
	ImageURL                string       `gorm:"column:image_url;default:default-member-profile.jpeg" json:"imageURL"`
	Description             string       `gorm:"column:description" json:"description"`
	UserId                  uuid.UUID    `gorm:"type:uuid;column:user_id" json:"user_id"`
	LastTimeScheduleChanged time.Time    `gorm:"column:schedule_changed_at" json:"schedule_changed_at"`
	LastTimeGeneratedTask   *time.Time   `gorm:"column:generate_task_updated_at" json:"generate_task_updated_at"`
	Members                 []*Members   `gorm:"foreignKey:calendar_id" json:"members"`
	Leaves                  []*Leaves    `gorm:"foreignKey:calendar_id" json:"leaves"`
	Schedules               []*Schedules `gorm:"foreignKey:calendar_id" json:"schedules"`
	Tasks                   []*Tasks     `gorm:"foreignKey:calendar_id" json:"tasks"`
	BaseModel
}

func (cal *Calendars) BeforeCreate(tx *gorm.DB) (err error) {
	if cal.Id == uuid.Nil {
		cal.Id = uuid.New()
	}

	return nil
}

func (calendar *Calendars) AfterSave(tx *gorm.DB) error {
	if tx.Statement.Changed("schedule_changed_at") {
		logger.Message("RecurrentScheduleChanged", zap.String("calendarId", calendar.Id.String()))
	}

	return nil
}
