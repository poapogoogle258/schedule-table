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
		logger.Debug("migrating init member")

		data := readFile("./data/members.json")
		var members []*dao.Members
		if err := json.Unmarshal(data, &members); err != nil {
			return err
		}

		return tx.Select("Id", "CalendarId", "ImageURL", "Name", "Nickname", "Color", "Description", "Position", "Email", "Telephone").Create(&members).Error
	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollicking init member")

		data := readFile("./data/members.json")
		var members []*dao.Members
		if err := json.Unmarshal(data, &members); err != nil {
			return err
		}

		ids := make([]string, len(members))
		for i := range members {
			ids[i] = members[i].Id.String()
		}

		return tx.Unscoped().Delete(&dao.Members{}, "id IN ?", ids).Error

	},
}
