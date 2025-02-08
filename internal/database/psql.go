package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresql() *gorm.DB {

	var (
		name     = os.Getenv("DB_NAME")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", host, username, password, name, port)
	)

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if os.Getenv("MIGRATE_SETUP") == "init" {
		migrate_err := MigrateSetUpAndInitData(db)

		if migrate_err != nil {
			fmt.Println(migrate_err)
		}
	}

	return db

}
