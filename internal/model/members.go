package model

import "gorm.io/gorm"

type Members struct {
	gorm.Model
	Id          uint   `gorm:"column:id; primary_key; not null" json:"id"`
	CalendarId  uint   `gorm:"column:calendar_id" json:"calendar_id"`
	ImageURL    string `gorm:"column:imageURL" json:"imageURL"`
	Name        string `gorm:"column:name" json:"name"`
	Nickname    string `gorm:"column:nickname" json:"nickname"`
	Color       string `gorm:"column:color" json:"color"`
	Description string `gorm:"column:description" json:"description"`
	Position    string `gorm:"column:position" json:"position"`
	Email       string `gorm:"column:email" json:"email"`
	Telephone   string `gorm:"column:telephone" json:"telephone"`
	BaseModel
}
