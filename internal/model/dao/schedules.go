package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedules struct {
	Id                   uuid.UUID      `gorm:"type:uuid;column:id;primarykey;uniqueIndex" json:"id"`
	MasterScheduleId     *uuid.UUID     `gorm:"type:uuid;column:master_id" json:"master_id"`
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
	BreakTime            uint32         `gorm:"column:breaktime;default:0" json:"breaktime"`
	UseNumberPeople      int64          `gorm:"column:use_number_people;default:1" json:"use_number_people"`
	Recurrence_freq      int8           `gorm:"column:recurrence_freq" json:"recurrence_freq"` // YEARLY=0,MONTHLY,WEEKLY,DAILY,HOURLY,MINUTELY,SECONDLY
	Recurrence_interval  int32          `gorm:"column:recurrence_interval" json:"recurrence_interval"`
	Recurrence_count     int32          `gorm:"column:recurrence_count" json:"recurrence_count"`
	Recurrence_bymonth   string         `gorm:"column:recurrence_bymonth" json:"recurrence_bymonth"`
	Recurrence_byweekday string         `gorm:"column:recurrence_byweekday" json:"recurrence_byweekday"`
	Responsibles         *[]Responsible `gorm:"foreignKey:schedule_id;" json:"Responsibles"`
	CreatedAt            time.Time      `gorm:"index;column:created_at" json:"-"`
	UpdatedAt            time.Time      `gorm:"index;column:updated_at" json:"-"`
}

func (schedule *Schedules) BeforeCreate(tx *gorm.DB) (err error) {
	schedule.Id = uuid.New()

	return nil
}

func (schedule *Schedules) BeforeSave(tx *gorm.DB) (err error) {
	if schedule.Responsibles != nil {
		tx.Delete(&Responsible{}, "schedule_id = ?", schedule.Id)
	}

	return nil
}

func (schedule *Schedules) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Delete(&Responsible{}, "schedule_id = ?", schedule.Id)

	return nil
}
