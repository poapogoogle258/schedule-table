package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Leaves struct {
	gorm.Model
	Id          uuid.UUID `gorm:"type:uuid; column:id; primary_key; uniqueIndex" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	MemberId    uuid.UUID `gorm:"type:uuid;column:member_id" json:"member_id"`
	UserId      uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	Start       time.Time `gorm:"column:start" json:"start"`
	End         time.Time `gorm:"column:end" json:"end"`
	Description *string   `gorm:"column:description" json:"description"`
	Status      int8      `gorm:"column:status;default:0" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Member      Members   `gorm:"references:id"`
}
