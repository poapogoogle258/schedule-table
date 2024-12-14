package services

import (
	"schedule_table/internal/repository"

	"github.com/gin-gonic/gin"
)

type UsersService interface {
	GetAllUsers(c *gin.Context)
	GetUserId(c *gin.Context)
	UpdateUserId(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UsersServiceImpl struct {
	u repository.UserRepository
}

type RegisterAccountRequest struct {
}

func (u UsersServiceImpl) RegisterAccount(c *gin.Context) {

}
