package pkg

import (
	"fmt"
	"net/http"

	"schedule_table/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Null() interface{} {
	return nil
}

func PanicHandler(c *gin.Context) {
	if err := recover(); err != nil {
		// Get stack trace

		// Log console-friendly error (without stack)
		logger.Error(fmt.Sprintf("panic recovered: %v", err),
			zap.String("path", c.Request.URL.Path),
		)

		// Return error response
		msg := fmt.Sprintf("%s: %s", http.StatusText(500), err)
		c.JSON(http.StatusInternalServerError, BuildWithoutResponse(500, msg))
		c.Abort()
	}
}
