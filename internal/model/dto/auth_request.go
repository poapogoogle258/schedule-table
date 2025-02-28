package dto

import "github.com/google/uuid"

// signUpUserResponse represents the response payload for the sign-up endpoint.
type SignUpResponse struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ImageURL string    `json:"image"`
}

type SignUpRequest struct {
	Name        string `json:"name" binding:"required"`
	ImageUrl    string `json:"imageUrl" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Description string `json:"description"`
}

type LoginResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ImageURL    string    `json:"image"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   int64     `json:"expires_at"`
}

// profileResponse represents the response payload for the getProfile endpoint.
type UserProfile struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ImageURL    string `json:"imageURL"`
	Description string `json:"description"`
	CalendarId  string `json:"calendar_id"`
}
