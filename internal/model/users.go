package model

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Id          uint      `gorm:"column:id; primary_key; not null" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	ImageURL    string    `gorm:"column:imageURL; default:default-member-profile.jpeg" json:"imageURL"`
	Description string    `gorm:"column:description" json:"description"`
	Email       string    `gorm:"column:email" json:"email"`
	Telephone   string    `gorm:"column:telephone" json:"telephone"`
	Token       string    `gorm:"column:token" json:"token"`
	Calendar    Calendars `gorm:"foreignKey:user_id" json:"calendar"`
	BaseModel
}
