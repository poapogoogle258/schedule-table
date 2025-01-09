package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Leaves struct {
	Id          uuid.UUID `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	CalendarId  uuid.UUID `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	MemberId    uuid.UUID `gorm:"type:uuid;column:member_id" json:"member_id"`
	UserId      uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	Date        time.Time `gorm:"column:date" json:"date"`
	Tzid        string    `gorm:"column:tzid;default:Asia/Bangkok" json:"tzid"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Member      Members   `gorm:"references:id"`
}

func (leave *Leaves) BeforeCreate(tx *gorm.DB) (err error) {
	leave.Id = uuid.New()

	return
}
