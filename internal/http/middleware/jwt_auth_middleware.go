package middleware

import (
	"net/http"
	"schedule_table/internal/handler"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(handlerAuth handler.AuthHandler) gin.HandlerFunc {

	jwt_service := service.NewJWTAuthService()

	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := jwt_service.ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(*service.AuthCustomClaims)

			if err := handlerAuth.CheckUserTokenExist(claims, tokenString); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"statusCode": http.StatusUnauthorized,
					"message":    err.Error(),
				})
				c.Abort()
			} else {
				c.Set("requestAuthUserId", claims.UserId)
				c.Next()
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"message":    err.Error(),
			})
			c.Abort()
		}
	}
}
