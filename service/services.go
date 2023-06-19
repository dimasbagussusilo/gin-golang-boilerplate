package service

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"gorm.io/gorm"
)

type Services struct {
	UserService  IRepository
	TokenService IRepository
}

func Init(db *gorm.DB) *Services {
	return &Services{
		UserService:  NewRepository(&models.User{}, db),
		TokenService: NewRepository(&models.Token{}, db),
	}
}
