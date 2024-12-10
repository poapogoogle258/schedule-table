package model

import "time"

type Responsible struct {
	MemberId     uint      `gorm:"primaryKey" json:"member_id"`
	ScheduleId   uint      `gorm:"primaryKey" json:"schedule_id"`
	Queue        int8      `gorm:"column:queue" json:"queue"`
	Limit        int8      `gorm:"column:limit" json:"limit"`
	LastTimeTask time.Time `gorm:"column:lastTimeTask" json:"lastTimeTask"`
}
