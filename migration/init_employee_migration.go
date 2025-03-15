package main

import (
	"encoding/json"
	"schedule_table/internal/model/dao"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var MemberInit = &gormigrate.Migration{
	ID: "000000000003",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating init employee")

		data := readFile("./data/employees.json")
		var employees []*dao.Employee
		if err := json.Unmarshal(data, &employees); err != nil {
			return err
		}

		return tx.Select("Id", "CalendarId", "ImageURL", "Name", "Nickname", "Description", "Position", "Email", "Telephone").Create(&employees).Error
	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollicking init employee")

		data := readFile("./data/employees.json")
		var employees []*dao.Employee
		if err := json.Unmarshal(data, &employees); err != nil {
			return err
		}

		ids := make([]string, len(employees))
		for i := range employees {
			ids[i] = employees[i].Id.String()
		}

		return tx.Unscoped().Delete(&dao.Employee{}, "id IN ?", ids).Error

	},
}
