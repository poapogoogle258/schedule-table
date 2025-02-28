package dao

import (
	"github.com/google/uuid"
)

type Calendar struct {
	Id          uuid.UUID   `gorm:"primarykey;type:uuid;column:id" json:"id"`
	Name        string      `gorm:"column:name" json:"name"`
	ImageURL    string      `gorm:"column:image_url;default:default-member-profile.jpeg" json:"imageURL"`
	Description string      `gorm:"column:description" json:"description"`
	UserId      uuid.UUID   `gorm:"type:uuid;column:user_id" json:"user_id"`
	Employees   []*Employee `gorm:"foreignKey:calendar_id" json:"members"`
	Leaves      []*Leave    `gorm:"foreignKey:calendar_id" json:"leaves"`
	Schedules   []*Schedule `gorm:"foreignKey:calendar_id" json:"schedules"`
	Tasks       []*Task     `gorm:"foreignKey:calendar_id" json:"tasks"`
	BaseModel
	SoftDelete
}

// func (calendar *Calendar) AfterSave(tx *gorm.DB) error {
// 	if tx.Statement.Changed("schedule_changed_at") {
// 		logger.Message("RecurrentScheduleChanged", zap.String("calendarId", calendar.Id.String()))
// 	}

// 	return nil
// }
