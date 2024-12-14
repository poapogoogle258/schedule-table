package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Calendars struct {
	gorm.Model
	Id          uuid.UUID    `gorm:"type:uuid; column:id; primary_key; uniqueIndex" json:"id"`
	Name        string       `gorm:"column:name" json:"name"`
	ImageURL    string       `gorm:"column:imageURL; default:default-member-profile.jpeg" json:"imageURL"`
	Description *string      `gorm:"column:description" json:"description"`
	UserId      uuid.UUID    `gorm:"type:uuid;column:user_id" json:"user_id"`
	Members     []*Members   `gorm:"foreignKey:calendar_id" json:"members"`
	Leaves      []*Leaves    `gorm:"foreignKey:calendar_id" json:"leaves"`
	Schedules   []*Schedules `gorm:"foreignKey:calendar_id" json:"schedules"`
	Tasks       []*Tasks     `gorm:"foreignKey:calendar_id" json:"tasks"`
	BaseModel
}
