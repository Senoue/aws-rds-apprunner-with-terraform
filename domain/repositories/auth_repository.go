package repository

import (
	"github.com/Senoue/aws-rds-apprunner-with-terraform/domain/models"
	"github.com/gin-gonic/gin"
)

// AuthRepository is an interface for authentication-related operations.
type AuthRepository interface {
	Login(ctx *gin.Context, email, password string) (*model.User, error)
	Register(ctx *gin.Context, user *model.User) error
	UserInfo(ctx *gin.Context, id int) (*model.User, error)
}
