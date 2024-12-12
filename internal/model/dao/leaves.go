package dao

import "time"

type Leaves struct {
	Id          uint      `gorm:"column:id; primary_key; not null" json:"id"`
	CalendarId  uint      `gorm:"column:calendar_id" json:"calendar_id"`
	MemberId    uint      `gorm:"column:member_id" json:"member_id"`
	UserId      uint      `gorm:"column:user_id" json:"user_id"`
	Start       time.Time `gorm:"column:start" json:"start"`
	End         time.Time `gorm:"column:end" json:"end"`
	Description string    `gorm:"column:description" json:"description"`
	Status      int8      `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Member      Members   `gorm:"references:id"`
}
