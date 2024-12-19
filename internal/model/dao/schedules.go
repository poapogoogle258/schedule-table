package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedules struct {
	gorm.Model
	Id                   uuid.UUID      `gorm:"type:uuid;column:id;primary_key;uniqueIndex" json:"id"`
	CalendarId           uuid.UUID      `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	Name                 string         `gorm:"column:name" json:"name"`
	Description          string         `gorm:"column:description;default:-" json:"description"`
	ImageURL             string         `gorm:"column:imageURL;default:default-image-schedule.jpeg" json:"imageURL"`
	Priority             int8           `gorm:"column:priority" json:"priority"`
	Start                time.Time      `gorm:"column:start" json:"start"`
	End                  time.Time      `gorm:"column:end" json:"end"`
	Hr_start             string         `gorm:"column:hr_start" json:"hr_start"`
	Hr_end               string         `gorm:"column:hr_end" json:"hr_end"`
	Tzid                 string         `gorm:"column:tzid;default:Asia/Bangkok" json:"tzid"`
	RestTime             int32          `gorm:"column:restTime;default:0" json:"restTime"`
	Recurrence_freq      int8           `gorm:"column:recurrence_freq" json:"recurrence_freq"` // YEARLY=0,MONTHLY,WEEKLY,DAILY,HOURLY,MINUTELY,SECONDLY
	Recurrence_interval  int32          `gorm:"column:recurrence_interval" json:"recurrence_interval"`
	Recurrence_wkst      string         `gorm:"column:recurrence_wkst" json:"recurrence_wkst"`
	Recurrence_bymonth   string         `gorm:"column:recurrence_bymonth" json:"recurrence_bymonth"`
	Recurrence_byweekday string         `gorm:"column:recurrence_byweekday" json:"recurrence_byweekday"`
	Members_responsible  *[]Responsible `gorm:"foreignKey:schedule_id;" json:"members_responsible"`
}
