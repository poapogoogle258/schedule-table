package model

import "time"

type Tasks struct {
	Id          uint      `gorm:"column:id; primary_key; not null" json:"id"`
	CalendarId  uint      `gorm:"column:calendar_id" json:"calendar_id"`
	ScheduleId  uint      `gorm:"column:schedule_id" json:"schedule_id"`
	MemberId    uint      `gorm:"column:member_id" json:"member_id"`
	Start       time.Time `gorm:"column:start" json:"start"`
	End         time.Time `gorm:"column:end" json:"end"`
	Status      int8      `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Person      Members   `gorm:"foreignKey:MemberId"`
	Description Schedules `gorm:"foreignKey:ScheduleId"`
}
