package database

import (
	"fmt"
	"os"
	"schedule_table/internal/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func mustConnectDb(dialector gorm.Dialector, opts ...gorm.Option) *gorm.DB {
	result, err := gorm.Open(dialector, opts...)
	if err != nil {
		logger.Error("databaseConnectFailed", zap.Error(err))
		panic(err)
	}

	return result
}

func ConnectPostgresql() *gorm.DB {

	if db == nil {
		var (
			name     = os.Getenv("DB_NAME")
			host     = os.Getenv("DB_HOST")
			port     = os.Getenv("DB_PORT")
			username = os.Getenv("DB_USERNAME")
			password = os.Getenv("DB_PASSWORD")
			dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", host, username, password, name, port)
		)

		db = mustConnectDb(postgres.Open(dsn), &gorm.Config{})

	}

	return db

}
