package forms

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
)

type WhoAmIResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type GetAllUserResponse struct {
	Users      []models.User     `json:"users"`
	Pagination *utils.Pagination `json:"pagination"`
}
