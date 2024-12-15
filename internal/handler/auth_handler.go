package handler

import (
	"net/http"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"schedule_table/unities"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type authHandler struct {
	jwtService service.JWTService
	userRepo   repository.UserRepository
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *authHandler) Login(c *gin.Context) {
	var request loginRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no data found",
		})

	}

	user := s.userRepo.FineUserEmail(request.Email)

	if unities.VerifyPassword(request.Password, user.Password) {
		token := s.jwtService.GenerateToken(user.Id.String(), user.Name, user.Email)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}

func (s *authHandler) ValidateToken(c *gin.Context) {

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

func NewAuthHandler(jwtService service.JWTService, userRepo repository.UserRepository) AuthHandler {
	return &authHandler{
		jwtService,
		userRepo,
	}
}
