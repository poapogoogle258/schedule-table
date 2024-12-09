package main

import (
	"fmt"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	addr := fmt.Sprint(os.Getenv("IP"), ":", os.Getenv("PORT"))

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 200,
			"message":    "success",
			"data":       "welcome to schedule tables project.",
		})
	})

	fmt.Println(addr)

	r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
