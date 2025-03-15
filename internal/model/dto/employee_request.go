package dto

type EmployeeInfoRequest struct {
	Name        string `json:"name" binding:"required"`
	NickName    string `json:"nickname" binding:"required"`
	Description string `json:"description" binding:"required"`
	Position    string `json:"position" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Telephone   string `json:"telephone" binding:"required"`
	ImageURL    string `json:"imageURL" binding:"required"`
}

type EmployeeInfo struct {
	Id string `json:"id"`
	EmployeeInfoRequest
}
