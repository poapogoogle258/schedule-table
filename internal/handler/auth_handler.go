package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"schedule_table/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var ErrRequestAuthorizationHeader = errors.New("request token in authorization header")

type AuthHandler interface {
	Login(c *gin.Context)
	ValidateToken(c *gin.Context)
	Profile(c *gin.Context)
	SignUp(c *gin.Context)
}

type AuthHandlerImpl struct {
	jwtService service.JWTService
	userRepo   repository.UserRepository
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpBody struct {
	Name        string `json:"name" binding:"required"`
	ImageUrl    string `json:"imageUrl" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Description string `json:"description"`
}

type SignUpUserResponse struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Image string    `json:"image"`
}

func (handler *AuthHandlerImpl) SignUp(c *gin.Context) {
	var body SignUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, pkg.BuildWithoutResponse(http.StatusBadRequest, err.Error()))
		c.Abort()

		return
	}

	if !handler.userRepo.IsUniqueEmail(body.Email) {
		c.JSON(http.StatusBadRequest, pkg.BuildWithoutResponse(http.StatusBadRequest, repository.ErrDuplicateEmail.Error()))
		c.Abort()

		return
	}

	newUser := &dao.Users{
		Name:        body.Name,
		ImageURL:    body.ImageUrl,
		Email:       body.Email,
		Password:    body.Password,
		Description: body.Description,
	}

	if err := handler.userRepo.Register(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, pkg.BuildWithoutResponse(http.StatusInternalServerError, err.Error()))
		c.Abort()

		return
	}

	if _, err := handler.userRepo.CreateCalendarDefault(newUser.Id); err != nil {
		c.JSON(http.StatusInternalServerError, pkg.BuildWithoutResponse(http.StatusInternalServerError, err.Error()))
		c.Abort()

		return
	}

	resp := &SignUpUserResponse{}
	copier.Copy(&resp, newUser)

	c.JSON(http.StatusOK, pkg.BuildResponse(http.StatusOK, resp))
}

func (handler *AuthHandlerImpl) Login(c *gin.Context) {

	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no data found",
		})
		return
	}

	user, err := handler.userRepo.FindOneByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
	}

	if util.VerifyPassword(request.Password, user.Password) {
		token := handler.jwtService.GenerateToken(user.Id.String(), user.Name, user.Email)
		if err := handler.userRepo.UpdateOne(user.Id.String(), "token", token); err != nil {
			panic(err)
		}

		decode, _ := handler.jwtService.ValidateToken(token)
		profile, _ := handler.userRepo.GetProfile(user.Id.String())

		c.JSON(http.StatusOK, gin.H{
			"id":           profile.Id,
			"name":         profile.Name,
			"email":        profile.Email,
			"image":        profile.ImageURL,
			"access_token": token,
			"expires_at":   decode.Claims.(*service.AuthCustomClaims).ExpiresAt,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}

func (handler *AuthHandlerImpl) Profile(c *gin.Context) {

	tokenString, errGetToken := getTokenFromHeader(c)
	if errGetToken != nil {
		c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, errGetToken.Error()))
		c.Abort()

		return
	}

	if token, err := handler.jwtService.ValidateToken(tokenString); err == nil {
		claims := token.Claims.(*service.AuthCustomClaims)

		if profile, err := handler.userRepo.GetProfile(claims.UserId); err != nil {
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

	tokenString, errGetToken := getTokenFromHeader(c)
	if errGetToken != nil {
		c.JSON(http.StatusUnauthorized, pkg.BuildWithoutResponse(http.StatusUnauthorized, errGetToken.Error()))
		c.Abort()

		return
	}

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

func getTokenFromHeader(c *gin.Context) (string, error) {

	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || len(authHeader) <= len(BEARER_SCHEMA) {
		return "", ErrRequestAuthorizationHeader
	}

	return authHeader[len(BEARER_SCHEMA):], nil
}

func NewAuthHandler(jwtService service.JWTService, userRepo repository.UserRepository) AuthHandler {
	return &AuthHandlerImpl{
		jwtService,
		userRepo,
	}
}
