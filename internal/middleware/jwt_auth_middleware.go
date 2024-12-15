package middleware

import (
	"net/http"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := service.NewJWTAuthService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(*service.AuthCustomClaims)

			c.Keys = make(map[string]any)
			c.Keys["token_userId"] = claims.UserId
			c.Keys["token_name"] = claims.Name
			c.Keys["token_email"] = claims.Email
			c.Next()

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"message":    err.Error(),
			})

		}
	}
}
