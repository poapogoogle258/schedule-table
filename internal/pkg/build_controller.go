package pkg

import (
	"net/http"
	"schedule_table/internal/constant"

	"github.com/gin-gonic/gin"
)

func BuildController(method string, handler func(c *gin.Context) (interface{}, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			if method == "Post" {
				c.JSON(http.StatusCreated, BuildResponse(constant.Success, response))
			} else if method == "Delete" {
				c.JSON(http.StatusNoContent, BuildResponse(constant.Success, 0))
			} else {
				c.JSON(http.StatusOK, BuildResponse(constant.Success, response))
			}
		}
	}
}

func BuildGetController(handler func(c *gin.Context) (interface{}, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusOK, BuildResponse(constant.Success, response))
		}
	}
}

func BuildPostController(handler func(c *gin.Context) (interface{}, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer PanicHandler(c)

		if response, err := handler(c); err != nil {
			panic(err)
		} else {
			c.JSON(http.StatusCreated, BuildResponse(constant.Success, response))
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
