package model

import "gorm.io/gorm"

type Calendars struct {
	gorm.Model
	Id          uint       `gorm:"column:id; primary_key; not null" json:"id"`
	Name        string     `gorm:"column:name" json:"name"`
	ImageURL    string     `gorm:"column:imageURL; default:default-member-profile.jpeg" json:"imageURL"`
	Description string     `gorm:"column:description" json:"description"`
	UserId      uint       `gorm:"column:user_id"`
	Owner       Users      `gorm:"references:UserId" json:"owner"`
	Members     []*Members `gorm:"polymorphic:Members" json:"members"`
	BaseModel
}
