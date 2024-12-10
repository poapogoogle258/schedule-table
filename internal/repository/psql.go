package repository

import (
	"fmt"
	"os"

	"github.com/poapogoogle258/schedule_table/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresql() (*gorm.DB, error) {

	var (
		name     = os.Getenv("DB_NAME")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", host, username, password, name, port)
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	migrate_err := db.AutoMigrate(&model.Users{}, &model.Leaves{}, &model.Members{}, &model.Schedules{}, &model.Responsible{}, &model.Tasks{}, &model.Calendars{})

	if migrate_err != nil {
		fmt.Println(migrate_err)
	}

	return db, err

}
