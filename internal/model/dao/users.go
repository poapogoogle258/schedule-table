package dao

import (
	"github.com/google/uuid"
)

type Users struct {
	Id          uuid.UUID  `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	Name        string     `gorm:"column:name" json:"name"`
	ImageURL    string     `gorm:"column:imageURL;default:default-member-profile.jpeg" json:"imageURL"`
	Description string     `gorm:"column:description" json:"description"`
	Email       string     `gorm:"column:email;not null;uniqueIndex" json:"email"`
	Password    string     `gorm:"column:password" json:"-"`
	Telephone   string     `gorm:"column:telephone" json:"telephone"`
	Token       string     `gorm:"column:token" json:"token"`
	Calendar    *Calendars `gorm:"foreignKey:user_id" json:"calendar"`
	BaseModel
}
