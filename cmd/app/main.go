package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/poapogoogle258/schedule_table/internal/repository"
)

func main() {
	godotenv.Load()

	// r := gin.Default()
	// addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))

	_, err := repository.ConnectPostgresql()

	if err != nil {
		fmt.Println(err)
	}

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"statusCode": 200,
	// 		"message":    "success",
	// 		"data":       "welcome to schedule tables project.",
	// 	})
	// })

	// r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
