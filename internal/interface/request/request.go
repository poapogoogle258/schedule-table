package request

import (
	"fmt"
	"mime/multipart"
	"regexp"
)

type Member struct {
	ImageURL    string `json:"imageURL"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Telephone   string `json:"telephone"`
}

type CreateNewMember struct {
	Name        string                `form:"name"`
	NickName    string                `form:"nickname"`
	Color       string                `form:"color"`
	Description string                `form:"description"`
	Position    string                `form:"position"`
	Email       string                `form:"email"`
	Telephone   string                `form:"telephone"`
	File        *multipart.FileHeader `form:"image"`
}

func (newMember *CreateNewMember) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if newMember.Email != "" && !emailRegex.MatchString(newMember.Email) {
		return fmt.Errorf("Validate: field 'Email' format validate failed")
	}

	colorRegex := regexp.MustCompile(`^#[a-fA-F0-9]{6}$`)
	if newMember.Color != "" && !colorRegex.MatchString(newMember.Color) {
		return fmt.Errorf("Validate: field 'Color' format validate failed")
	}

	telephoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if newMember.Telephone != "" && !telephoneRegex.MatchString(newMember.Telephone) {
		return fmt.Errorf("Validate: field 'Telephone' format validate failed")
	}

	return nil
}
