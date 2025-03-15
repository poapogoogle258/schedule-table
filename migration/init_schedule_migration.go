package main

import (
	"encoding/json"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/util"
	"strconv"
	"strings"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var scheduleColumns = []string{"Id", "MasterScheduleId", "CalendarId", "Name", "Description", "ImageURL", "Priority", "Start", "End", "Hr_start", "Hr_end", "Tzid", "BreakTime", "UseNumberPeople", "Recurrence_freq", "Recurrence_interval", "Recurrence_count", "Recurrence_bymonth", "Recurrence_byweekday", "Color", "RotationCycle"}

type member struct {
	Id          string `json:"id"`
	ImageURL    string `json:"imageURL"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Telephone   string `json:"telephone"`
}

type recurrence struct {
	Freq      int8  `json:"freq"`
	Count     int32 `json:"count"`
	Interval  int32 `json:"interval"`
	Byweekday []int `json:"byweekday"`
	Bymonth   []int `json:"bymonth"`
}

type Schedule struct {
	Id               string     `json:"id"`
	MasterScheduleId string     `json:"master_id"`
	CalendarId       string     `json:"calendar_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	ImageURL         string     `json:"imageURL"`
	Priority         int8       `json:"priority"`
	Start            time.Time  `json:"start"`
	End              *time.Time `json:"end"`
	Hr_start         string     `json:"hr_start"`
	Hr_end           string     `json:"hr_end"`
	Tzid             string     `json:"tzid"`
	BreakTime        uint32     `json:"breaktime"`
	UseNumberPeople  int8       `json:"use_number_people"`
	Recurrence       recurrence `json:"recurrence"`
	Members          []member   `json:"members"`
	RotationCycle    int        `json:"rotation_cycle"`
	Color            string     `json:"color"`
}

// json -> dao
func (s Schedule) Recurrence_freq() int8 {
	return s.Recurrence.Freq
}
func (s Schedule) Recurrence_interval() int32 {
	return s.Recurrence.Interval
}
func (s Schedule) Recurrence_count() int32 {
	return s.Recurrence.Count
}
func (s Schedule) Recurrence_bymonth() string {
	if len(s.Recurrence.Bymonth) == 0 {
		return ""
	}
	strList := util.Map(s.Recurrence.Bymonth, func(i int) string {
		return strconv.Itoa(i)
	})

	return strings.Join(strList, ",")
}
func (s Schedule) Recurrence_byweekday() string {
	if len(s.Recurrence.Byweekday) == 0 {
		return ""
	}
	strList := util.Map(s.Recurrence.Byweekday, func(i int) string {
		return strconv.Itoa(i)
	})

	return strings.Join(strList, ",")
}

var ScheduleInit = &gormigrate.Migration{
	ID: "000000000004",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating init schedule")

		data := readFile("./data/schedules.json")
		var schedules []*Schedule
		if err := json.Unmarshal(data, &schedules); err != nil {
			return err
		}
		insert := []*dao.Schedule{}
		if err := copier.Copy(&insert, schedules); err != nil {
			return err
		}

		// insert schedule

		if err := tx.Select(scheduleColumns).Create(&insert).Error; err != nil {
			tx.Rollback()
			return err
		}

		// insert responsible
		for i := range schedules {
			if schedules[i].MasterScheduleId == schedules[i].Id {
				responsible := make([]*dao.EmployeeQueue, len(schedules[i].Members))
				for j := range responsible {
					responsible[j] = &dao.EmployeeQueue{
						ScheduleId: uuid.MustParse(schedules[i].Id),
						EmployeeId: uuid.MustParse(schedules[i].Members[j].Id),
						Queue:      int8(j),
					}
				}

				if err := tx.Create(&responsible).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		return nil

	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollicking init schedule")

		data := readFile("./data/schedules.json")
		var schedules []*dto.ScheduleInfo
		if err := json.Unmarshal(data, &schedules); err != nil {
			return err
		}
		insert := []*dao.Schedule{}
		if err := copier.Copy(&insert, schedules); err != nil {
			return err
		}

		ids := make([]string, len(insert))
		for i := range insert {
			ids[i] = insert[i].Id.String()
		}

		if err := tx.Where("schedule_id In ?", ids).Delete(&dao.EmployeeQueue{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&dao.Schedule{}, "id IN ?", ids).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil

	},
}
