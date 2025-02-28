package main

import (
	"fmt"
	"os"
	"schedule_table/internal/http/router"
	"schedule_table/internal/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")

	handlers := Injector()
	server := router.NewRouter(handlers)

	// initial logger
	logger.InitLogger()
	defer logger.Sync()

	// start http servers
	addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))
	server.Run(addr) // listen and serve on 0.0.0.0:8080

}
