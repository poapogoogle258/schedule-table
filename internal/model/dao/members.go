package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Members struct {
	Id           uuid.UUID  `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	CalendarId   uuid.UUID  `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	ImageURL     string     `gorm:"column:image_url;default: default-member-profile.jpeg" json:"imageURL"`
	Name         string     `gorm:"column:name" json:"name"`
	Nickname     string     `gorm:"column:nickname" json:"nickname"`
	Color        string     `gorm:"column:color;default:#000000" json:"color"`
	Description  string     `gorm:"column:description" json:"description"`
	Position     string     `gorm:"column:position" json:"position"`
	Email        string     `gorm:"column:email" json:"email"`
	Telephone    string     `gorm:"column:telephone" json:"telephone"`
	LastTimeTask *time.Time `gorm:"column:lastTimeTask" json:"lastTimeTask"`
	Leaves       *[]Leaves  `gorm:"foreignKey:member_id" json:"leaves"`
	BaseModel
}

func (mem *Members) BeforeCreate(tx *gorm.DB) (err error) {
	mem.Id = uuid.New()

	return
}
