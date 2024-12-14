package dao

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
