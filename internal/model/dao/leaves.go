package dao

import (
	"time"

	"github.com/google/uuid"
)

type Leaves struct {
	Id          uuid.UUID `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	MemberId    uuid.UUID `gorm:"type:uuid;column:member_id" json:"member_id"`
	UserId      uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	Start       time.Time `gorm:"column:start" json:"start"`
	End         time.Time `gorm:"column:end" json:"end"`
	Description string    `gorm:"column:description" json:"description"`
	Status      int8      `gorm:"column:status;default:0" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Member      Members   `gorm:"references:id"`
}
