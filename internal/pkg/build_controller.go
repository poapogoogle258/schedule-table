package pkg

import (
	"net/http"
	"schedule_table/internal/constant"

	"github.com/gin-gonic/gin"
)

func BuildGetController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusOK, BuildResponse(constant.Success, response))
		}
	}
}

func BuildPostController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusCreated, BuildResponse(constant.Success, response))
		}
	}
}

func BuildPatchController[T any](handler func(c *gin.Context) (T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusOK, BuildResponse(constant.Success, response))
		}
	}
}

func BuildDeleteController(handler func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusNoContent, BuildResponse(constant.Success, ""))
		}
	}
}
