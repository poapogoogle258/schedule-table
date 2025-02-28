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

		return tx.AutoMigrate(&dao.User{}, &dao.Leave{}, &dao.Employee{}, &dao.Schedule{}, &dao.EmployeeQueue{}, &dao.Task{}, &dao.Calendar{})

	},
	Rollback: func(tx *gorm.DB) error {
		logger.Debug("rollback drop table database")

		return tx.Migrator().DropTable(&dao.User{}, &dao.Leave{}, &dao.Employee{}, &dao.Schedule{}, &dao.EmployeeQueue{}, &dao.Task{}, &dao.Calendar{})
	},
}
