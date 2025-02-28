package dao

import (
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `gorm:"primarykey;type:uuid;column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	ImageURL    string    `gorm:"column:imageURL;default:default-member-profile.jpeg" json:"imageURL"`
	Description string    `gorm:"column:description" json:"description"`
	Email       string    `gorm:"column:email;not null;uniqueIndex" json:"email"`
	Password    string    `gorm:"column:password" json:"password"`
	Token       string    `gorm:"column:token" json:"token"`
	Calendar    *Calendar `gorm:"foreignKey:user_id" json:"calendar"`
	BaseModel
	SoftDelete
}
