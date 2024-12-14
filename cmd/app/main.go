package main

import (
	"fmt"
	"os"
	"schedule_table/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server := router.NewRouter(Injector())

	addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))

	server.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
