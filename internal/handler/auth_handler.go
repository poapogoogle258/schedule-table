package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dto"
	"schedule_table/internal/pkg"
	"schedule_table/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	ErrRequestAuthorizationHeader = errors.New("request token in authorization header")
	ErrAuthEmailInvalid           = errors.New("email invalid")
	ErrAuthPasswordInvalid        = errors.New("password invalid")
	ErrAuthTokenInvalid           = errors.New("token invalid")
)

type AuthHandler interface {
	SignUp(c *gin.Context) (*dto.SignUpResponse, error)
	Login(c *gin.Context) (*dto.LoginResponse, error)
	GetProfile(c *gin.Context) (*dto.UserProfile, error)
}

// authHandler handles authentication-related requests such as login, sign-up, and profile retrieval.
type authHandler struct {
	jwtService  service.JWTService
	userService service.UserService
}

func (auth *authHandler) SignUp(c *gin.Context) (*dto.SignUpResponse, error) {
	var body dto.SignUpRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if err := auth.userService.ValidateNewUser(body.Name, body.Email, body.Password, body.ImageUrl); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	return auth.userService.RegisterUser(&body)
}

// loginRequest represents the payload for the login endpoint.
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (auth *authHandler) Login(c *gin.Context) (*dto.LoginResponse, error) {

	var body loginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	user, err := auth.userService.Authentication(body.Email, body.Password)
	if err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	token := auth.jwtService.GenerateToken(user.Id.String(), user.Name, user.Email)
	decode, errDecodeToken := auth.jwtService.ValidateToken(token)

	if errDecodeToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errDecodeToken)
	}

	// back-door auth
	if err := auth.userService.UpdateToken(user.Id.String(), token); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	resp := dto.LoginResponse{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		ImageURL:    user.ImageURL,
		AccessToken: token,
		ExpiresAt:   decode.Claims.(*service.AuthCustomClaims).ExpiresAt,
	}

	return &resp, nil

}

func (handler *authHandler) GetProfile(c *gin.Context) (*dto.UserProfile, error) {

	tokenString, errGetToken := getTokenFromHeader(c)
	if errGetToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusUnauthorized, errGetToken)
	}

	token, errValidateToken := handler.jwtService.ValidateToken(tokenString)
	if errValidateToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusForbidden, ErrAuthTokenInvalid)
	}

	claims := token.Claims.(*service.AuthCustomClaims)
	return handler.userService.GetProfile(claims.UserId)

}

func NewAuthHandler(jwtService service.JWTService, userService service.UserService) AuthHandler {
	return &authHandler{
		jwtService,
		userService,
	}
}

func getTokenFromHeader(c *gin.Context) (string, error) {

	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	if len(authHeader) <= len(BEARER_SCHEMA) {
		return "", errors.New("authorization header is invalid")
	}

	return authHeader[len(BEARER_SCHEMA):], nil
}
