package dao

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"index;column:created_at" json:"-"`
	UpdatedAt time.Time      `gorm:"index;column:updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}
