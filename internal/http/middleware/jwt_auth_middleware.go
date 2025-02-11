package middleware

import (
	"errors"
	"net/http"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

type IAuthorizeJWTMiddleware interface {
	Authorize() gin.HandlerFunc
}

type AuthorizeJWTMiddleware struct {
	JwtService service.JWTService
	UserRepo   repository.UserRepository
}

var ErrTokenNotEqualUserToken = errors.New("token have changed, type login again")

func (auth *AuthorizeJWTMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]

		if token, err := auth.JwtService.ValidateToken(tokenString); !token.Valid {
			c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return

		} else {
			claims := token.Claims.(*service.AuthCustomClaims)

			if userToken, err := auth.UserRepo.GetTokenUser(claims.UserId); err != nil {
				c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, err.Error()))
				return

			} else {

				if userToken != tokenString {
					c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, ErrTokenNotEqualUserToken.Error()))
					c.Abort()
					return
				}

				c.Set("requestAuthUserId", claims.UserId)
				c.Next()
			}

		}
	}
}

func NewAuthorizeJWTMiddleware(jwtServer service.JWTService, userRepo repository.UserRepository) IAuthorizeJWTMiddleware {
	return &AuthorizeJWTMiddleware{
		JwtService: jwtServer,
		UserRepo:   userRepo,
	}
}
