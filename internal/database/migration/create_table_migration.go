package main

import (
	"schedule_table/internal/model/dao"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var CreateTable = &gormigrate.Migration{
	ID: "000000000000",
	Migrate: func(tx *gorm.DB) error {
		logger.Debug("migrating table database")

		return tx.AutoMigrate(&dao.Users{}, &dao.Leaves{}, &dao.Members{}, &dao.Schedules{}, &dao.Responsible{}, &dao.Tasks{}, &dao.Calendars{})

	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollback drop table database")

		return tx.Migrator().DropTable(&dao.Users{}, &dao.Leaves{}, &dao.Members{}, &dao.Schedules{}, &dao.Responsible{}, &dao.Tasks{}, &dao.Calendars{})
	},
}
