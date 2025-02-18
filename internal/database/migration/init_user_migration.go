package main

import (
	"encoding/json"
	"schedule_table/internal/model/dao"
	"schedule_table/util"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var UserInit = &gormigrate.Migration{
	ID: "000000000001",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating init user")

		data := readFile("./data/users.json")
		var users []*dao.Users
		if err := json.Unmarshal(data, &users); err != nil {
			return err
		}

		// hash password
		for i := range users {
			users[i].Password = util.HashPassword(users[i].Password)
		}

		return tx.Create(&users).Error
	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollicking init user")

		data := readFile("./data/users.json")
		var users []*dao.Users
		if err := json.Unmarshal(data, &users); err != nil {
			return err
		}

		ids := make([]string, len(users))
		for i := range users {
			ids[i] = users[i].Id.String()
		}

		return tx.Unscoped().Delete(&dao.Users{}, "id IN ?", ids).Error

	},
}
