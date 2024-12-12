package dao

import "time"

type Responsible struct {
	MemberId     uint      `gorm:"column:member_id; primaryKey" json:"member_id"`
	ScheduleId   uint      `gorm:"column:schedule_id; primaryKey" json:"schedule_id"`
	Queue        int8      `gorm:"column:queue" json:"queue"`
	Limit        int8      `gorm:"column:limit" json:"limit"`
	LastTimeTask time.Time `gorm:"column:lastTimeTask" json:"lastTimeTask"`
	Person       Members   `gorm:"foreignKey:member_id" json:"person"`
}
