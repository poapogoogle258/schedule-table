package main

import (
	"encoding/json"
	"schedule_table/internal/model/dao"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CalendarInit = &gormigrate.Migration{
	ID: "000000000002",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating init calendar")

		data := readFile("./data/calendars.json")
		var calendars []*dao.Calendars
		if err := json.Unmarshal(data, &calendars); err != nil {
			return err
		}

		return tx.Select("Id", "Name", "ImageURL", "Description", "UserId").Create(&calendars).Error
	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollicking init calendar")

		data := readFile("./data/calendars.json")
		var calendars []*dao.Calendars
		if err := json.Unmarshal(data, &calendars); err != nil {
			return err
		}

		ids := make([]string, len(calendars))
		for i := range calendars {
			ids[i] = calendars[i].Id.String()
		}

		return tx.Unscoped().Delete(&dao.Calendars{}, "id IN ?", ids).Error

	},
}
