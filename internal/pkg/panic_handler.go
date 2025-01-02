package pkg

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"net/http"
)

func Null() interface{} {
	return nil
}

func PanicHandler(c *gin.Context) {
	if err := recover(); err != nil {
		msg := fmt.Sprintf("%s: %s", http.StatusText(500), err.(error))
		c.JSON(http.StatusBadRequest, BuildWithoutResponse(500, msg))
		c.Abort()
	}
}
