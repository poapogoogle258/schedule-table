package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedule struct {
	Id                   uuid.UUID        `gorm:"primarykey;type:uuid;column:id" json:"id"`
	CalendarId           uuid.UUID        `gorm:"type:uuid;column:calendar_id" json:"calendar_id"`
	MasterScheduleId     uuid.UUID        `gorm:"type:uuid;column:master_id" json:"master_id"`
	Color                string           `gorm:"column:color;default:#000000" json:"color"`
	Name                 string           `gorm:"column:name" json:"name"`
	Description          string           `gorm:"column:description;default:-" json:"description"`
	ImageURL             string           `gorm:"column:image_url;default:default-image-schedule.jpeg" json:"image_url"`
	Priority             int8             `gorm:"column:priority" json:"priority"`
	Start                time.Time        `gorm:"column:start" json:"start"`
	End                  *time.Time       `gorm:"column:end" json:"end"`
	Tzid                 string           `gorm:"column:tzid;default:Asia/Bangkok" json:"tzid"`
	BreakTime            uint32           `gorm:"column:breaktime;default:0" json:"breaktime"`
	UseNumberPeople      int8             `gorm:"column:use_number_people;default:1" json:"use_number_people"`
	Hr_start             string           `gorm:"column:hr_start" json:"hr_start"`
	Hr_end               string           `gorm:"column:hr_end" json:"hr_end"`
	Recurrence_freq      int8             `gorm:"column:recurrence_freq" json:"recurrence_freq"` // YEARLY=0,MONTHLY,WEEKLY,DAILY,HOURLY,MINUTELY,SECONDLY
	Recurrence_interval  int32            `gorm:"column:recurrence_interval" json:"recurrence_interval"`
	Recurrence_count     int32            `gorm:"column:recurrence_count" json:"recurrence_count"`
	Recurrence_bymonth   string           `gorm:"column:recurrence_bymonth" json:"recurrence_bymonth"`
	Recurrence_byweekday string           `gorm:"column:recurrence_byweekday" json:"recurrence_byweekday"`
	RotationCycle        int              `gorm:"column:rotation_cycle" json:"rotation_cycle"`
	RecurrenceUpdatedAt  time.Time        `gorm:"column:recurrence_updated_at" json:"recurrence_updated_at"`
	RegenerateUpdatedAt  *time.Time       `gorm:"column:regenerate_updated_at" json:"regenerate_updated_at"`
	EmployeeQueue        []*EmployeeQueue `gorm:"foreignKey:schedule_id" json:"employee_queue"`
	BaseModel
}

func (schedule *Schedule) BeforeCreate(tx *gorm.DB) (err error) {

	schedule.RecurrenceUpdatedAt = time.Now()

	return nil
}

func (schedule *Schedule) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Delete(&EmployeeQueue{}, "schedule_id = ?", schedule.Id)
	tx.Delete(&Task{}, "schedule_id = ?", schedule.Id)

	return nil
}

func (schedule *Schedule) AfterUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Priority", "Start", "End", "Hr_start", "Hr_end", "Tzid", "BreakTime", "Recurrence_freq", "Recurrence_interval", "Recurrence_count", "Recurrence_bymonth", "Recurrence_byweekday", "Responsibles") {
		tx.Model(&Calendar{}).Where("id = ?", schedule.CalendarId).Update("schedule_changed_at", time.Now())
	}

	return nil
}
