package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildGetController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			ErrorHandler(c, err)
		} else {
			c.JSON(http.StatusOK, BuildResponse(http.StatusOK, response))
		}
	}
}

func BuildPostController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			ErrorHandler(c, err)
		} else {
			c.JSON(http.StatusCreated, BuildResponse(http.StatusCreated, response))
		}
	}
}

func BuildPatchController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			ErrorHandler(c, err)
		} else {
			c.JSON(http.StatusOK, BuildResponse(http.StatusOK, response))
		}
	}
}

func BuildDeleteController(handler func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if err := handler(c); err != nil {
			ErrorHandler(c, err)
		} else {
			c.JSON(http.StatusNoContent, BuildWithoutResponse(http.StatusNoContent, ""))
		}
	}
}
