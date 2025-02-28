package dao

import (
	"github.com/google/uuid"
)

type Employee struct {
	Id          uuid.UUID `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	Available   bool      `gorm:"column:available;default:true" json:"available"`
	ImageURL    string    `gorm:"column:image_url;default: default-member-profile.jpeg" json:"imageURL"`
	Name        string    `gorm:"column:name" json:"name"`
	Nickname    string    `gorm:"column:nickname" json:"nickname"`
	Description string    `gorm:"column:description" json:"description"`
	Position    string    `gorm:"column:position" json:"position"`
	Email       string    `gorm:"column:email" json:"email"`
	Telephone   string    `gorm:"column:telephone" json:"telephone"`
	Leaves      []*Leave  `gorm:"foreignKey:member_id" json:"leaves"`
	BaseModel
	SoftDelete
}
