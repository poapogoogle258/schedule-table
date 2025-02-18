package main

import (
	"os"
	"schedule_table/internal/database"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return data
}

func main() {
	godotenv.Load()

	db := database.ConnectPostgresql()
	logger = zap.Must(zap.NewDevelopment())

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{CreateTable, UserInit})

	m.RollbackTo(CreateTable.ID)

	// if err := m.Migrate(); err == nil {
	// 	// logger.Fatal("Migration failed: %v", zap.Error(err))
	// 	m.RollbackTo(CreateTable.ID)
	// }

	logger.Info("Migration did run successfully")

}
