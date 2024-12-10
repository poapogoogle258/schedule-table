package model

import "time"

type Schedules struct {
	Id                   uint           `gorm:"column:id; primary_key; not null" json:"id"`
	CalendarId           uint           `gorm:"column:calendar_id" json:"calendar_id"`
	Name                 string         `gorm:"column:name" json:"name"`
	Description          string         `gorm:"column:description" json:"description"`
	ImageURL             string         `gorm:"column:imageURL" json:"imageURL"`
	Priority             int8           `gorm:"column:priority" json:"priority"`
	Start                time.Time      `gorm:"column:start" json:"start"`
	End                  time.Time      `gorm:"column:end" json:"end"`
	Hr_start             string         `gorm:"column:hr_start" json:"hr_start"`
	Hr_end               string         `gorm:"column:hr_end" json:"hr_end"`
	Time_duration        int32          `gorm:"column:time_duration" json:"time_duration"`
	Recurrence_freq      string         `gorm:"column:recurrence_freq" json:"recurrence_freq"`
	Recurrence_interval  string         `gorm:"column:recurrence_interval" json:"recurrence_interval"`
	Recurrence_wkst      string         `gorm:"column:recurrence_wkst" json:"recurrence_wkst"`
	Recurrence_bymonth   string         `gorm:"column:recurrence_bymonth" json:"recurrence_bymonth"`
	Recurrence_byweekday string         `gorm:"column:recurrence_byweekday" json:"recurrence_byweekday"`
	Members_responsible  *[]Responsible `gorm:"foreignKey:schedule_id;" json:"members_responsible"`
}
