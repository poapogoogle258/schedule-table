package middleware

import (
	"errors"
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

type IAuthorizeJWTMiddleware interface {
	Authorize() gin.HandlerFunc
}

type AuthorizeJWTMiddleware struct {
	JwtService service.JWTService
}

var ErrTokenNotEqualUserToken = errors.New("token have changed, type login again")
var ErrRequestAuthorizationHeader = errors.New("request token in authorization header")

func (auth *AuthorizeJWTMiddleware) Authorize() gin.HandlerFunc {
	const BEARER_SCHEMA = "Bearer "

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) <= len(BEARER_SCHEMA) {
			c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, ErrRequestAuthorizationHeader.Error()))
			c.Abort()

			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]

		if token, err := auth.JwtService.ValidateToken(tokenString); !token.Valid {
			c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, err.Error()))
			c.Abort()

			return

		} else {
			claims := token.Claims.(*service.AuthCustomClaims)
			c.Set("authUserId", claims.UserId)
			c.Next()

		}
	}
}

func NewAuthorizeJWTMiddleware(jwtServer service.JWTService) IAuthorizeJWTMiddleware {
	return &AuthorizeJWTMiddleware{
		JwtService: jwtServer,
	}
}
