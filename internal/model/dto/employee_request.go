package dto

type EmployeeInfoRequest struct {
	Name        string `json:"name" validate:"required"`
	NickName    string `json:"nickname" validate:"required"`
	Description string `json:"description"`
	Position    string `json:"position" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Telephone   string `json:"telephone" validate:"required,telephone"`
	ImageURL    string `json:"imageURL" validate:"required,url"`
}

type EmployeeInfo struct {
	Id string `json:"id"`
	EmployeeInfoRequest
}
