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

	r := gin.Default()
	addr := fmt.Sprintf("%s:%s", os.Getenv("IP"), os.Getenv("PORT"))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 200,
			"message":    "success",
			"data":       "welcome to schedule tables project.",
		})
	})

	r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
