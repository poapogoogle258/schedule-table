package pkg

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandler(c *gin.Context, err error) {
	var (
		code int
		// message string
	)

	switch {
	case errors.As(err, &ErrorWithStatusCode{}):
		code = err.(ErrorWithStatusCode).Code
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = 400
	default:
		code = 500
	}

	c.JSON(code, BuildWithoutResponse(code, err.Error()))
	c.Abort()
}
