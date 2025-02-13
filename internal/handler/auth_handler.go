package handler

import (
	"errors"
	"net/http"
	"schedule_table/internal/model/dao"
	"schedule_table/internal/pkg"
	"schedule_table/internal/repository"
	"schedule_table/internal/service"

	"schedule_table/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var (
	ErrRequestAuthorizationHeader = errors.New("request token in authorization header")
	ErrAuthEmailInvalid           = errors.New("email invalid")
	ErrAuthPasswordInvalid        = errors.New("password invalid")
	ErrAuthTokenInvalid           = errors.New("token invalid")
)

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

// AuthHandler defines the methods related to user authentication, including login, profile retrieval, and sign-up.
type AuthHandler interface {
	SignUp(c *gin.Context) (*signUpUserResponse, error)
	Login(c *gin.Context) (*loginResponse, error)
	Profile(c *gin.Context) (*profileResponse, error)
}

// authHandler handles authentication-related requests such as login, sign-up, and profile retrieval.
type authHandler struct {
	jwtService service.JWTService
	userRepo   repository.UserRepository
}

// loginRequest represents the payload for the login endpoint.
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signUpBody represents the request body for the SignUp endpoint, containing user details required for registration.
type signUpBody struct {
	Name        string `json:"name" binding:"required"`
	ImageUrl    string `json:"imageUrl" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Description string `json:"description"`
}

// signUpUserResponse represents the response payload for the sign-up endpoint.
type signUpUserResponse struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL string    `json:"image"`
}

// loginResponse represents the response payload for the sign-in endpoint.
type loginResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ImageURL    string    `json:"image"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   int64     `json:"expires_at"`
}

// profileResponse represents the response payload for the getProfile endpoint.
type profileResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	ImageURL    string  `json:"imageURL"`
	Description string  `json:"description"`
	CalendarId  *string `json:"calendar_id"`
}

func (_profileResponse *profileResponse) Calendar(cal *dao.Calendars) {
	if cal != nil {
		id := cal.Id.String()
		_profileResponse.CalendarId = &id
	} else {
		_profileResponse.CalendarId = nil
	}
}

func (handler *authHandler) SignUp(c *gin.Context) (*signUpUserResponse, error) {
	var body signUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, err)
	}

	if !handler.userRepo.IsUniqueEmail(body.Email) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusBadRequest, repository.ErrDuplicateEmail)
	}

	newUser := &dao.Users{
		Id:          uuid.New(),
		Name:        body.Name,
		ImageURL:    body.ImageUrl,
		Email:       body.Email,
		Password:    util.HashPassword(body.Password),
		Description: body.Description,
	}

	if err := handler.userRepo.Create(newUser); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	if _, err := handler.userRepo.CreateCalendarDefault(newUser.Id); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	resp := signUpUserResponse{}
	copier.Copy(&resp, newUser)

	return &resp, nil
}

func (handler *authHandler) Login(c *gin.Context) (*loginResponse, error) {

	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	user, errQueryUserByEmail := handler.userRepo.FindOneByEmail(request.Email)
	if errQueryUserByEmail != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusForbidden, ErrAuthEmailInvalid)
	}

	if !util.VerifyPassword(request.Password, user.Password) {
		return nil, pkg.NewErrorWithStatusCode(http.StatusForbidden, ErrAuthPasswordInvalid)
	}

	token := handler.jwtService.GenerateToken(user.Id.String(), user.Name, user.Email)
	decode, errDecodeToken := handler.jwtService.ValidateToken(token)

	if errDecodeToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errDecodeToken)
	}
	profile, errGetProfile := handler.userRepo.GetProfile(user.Id.String())
	if errGetProfile != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, errGetProfile)
	}

	// back-door auth
	if err := handler.userRepo.UpdateOne(user.Id.String(), "token", token); err != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusInternalServerError, err)
	}

	resp := loginResponse{
		Id:          profile.Id,
		Name:        profile.Name,
		Email:       profile.Email,
		ImageURL:    profile.ImageURL,
		AccessToken: token,
		ExpiresAt:   decode.Claims.(*service.AuthCustomClaims).ExpiresAt,
	}

	return &resp, nil

}

func (handler *authHandler) Profile(c *gin.Context) (*profileResponse, error) {

	tokenString, errGetToken := getTokenFromHeader(c)
	if errGetToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusUnauthorized, errGetToken)
	}

	token, errValidateToken := handler.jwtService.ValidateToken(tokenString)
	if errValidateToken != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusForbidden, ErrAuthTokenInvalid)
	}

	claims := token.Claims.(*service.AuthCustomClaims)
	profile, errGetProfile := handler.userRepo.GetProfile(claims.UserId)
	if errGetProfile != nil {
		return nil, pkg.NewErrorWithStatusCode(http.StatusForbidden, errGetProfile)
	}

	resp := profileResponse{}
	copier.Copy(&resp, profile)

	return &resp, nil
}

func NewAuthHandler(jwtService service.JWTService, userRepo repository.UserRepository) AuthHandler {
	return &authHandler{
		jwtService,
		userRepo,
	}
}
