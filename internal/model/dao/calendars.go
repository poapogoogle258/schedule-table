package dao

import (
	"github.com/google/uuid"
)

type Calendars struct {
	Id          uuid.UUID    `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	Name        string       `gorm:"column:name" json:"name"`
	ImageURL    string       `gorm:"column:image_url;default:default-member-profile.jpeg" json:"imageURL"`
	Description string       `gorm:"column:description" json:"description"`
	UserId      uuid.UUID    `gorm:"type:uuid;column:user_id" json:"user_id"`
	Members     []*Members   `gorm:"foreignKey:calendar_id" json:"members"`
	Leaves      []*Leaves    `gorm:"foreignKey:calendar_id" json:"leaves"`
	Schedules   []*Schedules `gorm:"foreignKey:calendar_id" json:"schedules"`
	Tasks       []*Tasks     `gorm:"foreignKey:calendar_id" json:"tasks"`
	BaseModel
}
