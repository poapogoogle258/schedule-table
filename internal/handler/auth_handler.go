package handler

import (
	"fmt"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"schedule_table/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	ValidateToken(c *gin.Context)
	CheckUserTokenExist(claims *service.AuthCustomClaims, token string) error
	Profile(c *gin.Context)
}

type AuthHandlerImpl struct {
	jwtService service.JWTService
	userRepo   repository.UserRepository
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthHandlerImpl) Login(c *gin.Context) {

	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no data found",
		})
		return
	}

	user, err := s.userRepo.FindOneByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
	}

	if util.VerifyPassword(request.Password, user.Password) {
		token := s.jwtService.GenerateToken(user.Id.String(), user.Name, user.Email)
		if err := s.userRepo.UpdateOne(user.Id.String(), "token", token); err != nil {
			panic(err)
		}

		decode, _ := s.jwtService.ValidateToken(token)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"exp":   decode.Claims.(*service.AuthCustomClaims).ExpiresAt,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}

func (handler *AuthHandlerImpl) Profile(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]

	if token, err := handler.jwtService.ValidateToken(tokenString); err == nil {
		claims := token.Claims.(*service.AuthCustomClaims)

		if profile, err := handler.userRepo.Profile(claims.UserId); err != nil {
			c.JSON(http.StatusForbidden, pkg.BuildWithoutResponse(http.StatusForbidden, err.Error()))
			return
		} else {
			response := util.Convert[dto.ResponseProfile](&profile)
			c.JSON(http.StatusOK, pkg.BuildResponse(http.StatusOK, response))
			return
		}
	} else {
		c.JSON(http.StatusForbidden, pkg.BuildWithoutResponse(http.StatusForbidden, "token invalid"))
	}
}

func (s *AuthHandlerImpl) ValidateToken(c *gin.Context) {

	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]

	token, err := s.jwtService.ValidateToken(tokenString)

	if token.Valid {
		claims := token.Claims.(*service.AuthCustomClaims)

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"message":    "success",
			"data":       claims,
		})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    err.Error(),
		})

	}

}

func (s *AuthHandlerImpl) CheckUserTokenExist(claims *service.AuthCustomClaims, token string) error {
	user, err := s.userRepo.FindOne(claims.UserId)
	if err != nil {
		return fmt.Errorf("not found this user")
	} else if user.Token != token {
		return fmt.Errorf("token duplicate, try login again")
	} else {
		return nil
	}
}

func NewAuthHandler(jwtService service.JWTService, userRepo repository.UserRepository) AuthHandler {
	return &AuthHandlerImpl{
		jwtService,
		userRepo,
	}
}
