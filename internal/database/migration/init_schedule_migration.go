package main

import (
	"encoding/json"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var scheduleColumns = []string{"Id", "MasterScheduleId", "CalendarId", "Name", "Description", "ImageURL", "Priority", "Start", "End", "Hr_start", "Hr_end", "Tzid", "BreakTime", "UseNumberPeople", "Recurrence_freq", "Recurrence_interval", "Recurrence_count", "Recurrence_bymonth", "Recurrence_byweekday"}

var ScheduleInit = &gormigrate.Migration{
	ID: "000000000004",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating init schedule")

		data := readFile("./data/schedules.json")
		var schedules []*dto.ResponseSchedule
		if err := json.Unmarshal(data, &schedules); err != nil {
			return err
		}
		insert := []*dao.Schedules{}
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
			if schedules[i].MasterScheduleId == nil {
				responsible := make([]*dao.Responsible, len(schedules[i].Members))
				for j := range responsible {
					responsible[j] = &dao.Responsible{
						ScheduleId: uuid.MustParse(schedules[i].Id),
						MemberId:   uuid.MustParse(schedules[i].Members[j].Id),
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
		var schedules []*dto.ResponseSchedule
		if err := json.Unmarshal(data, &schedules); err != nil {
			return err
		}
		insert := []*dao.Schedules{}
		if err := copier.Copy(&insert, schedules); err != nil {
			return err
		}

		ids := make([]string, len(insert))
		for i := range insert {
			ids[i] = insert[i].Id.String()
		}

		if err := tx.Where("schedule_id In ?", ids).Delete(&dao.Responsible{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&dao.Schedules{}, "id IN ?", ids).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil

	},
}
