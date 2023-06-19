package controllers

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/forms"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	c  *config.Config
	db *gorm.DB
	s  *service.Services
}

func NewUserController(config *config.Config, db *gorm.DB, s *service.Services) *UserController {
	return &UserController{
		c:  config,
		db: db,
		s:  s,
	}
}

// Me godoc
// @Summary Show logged-in user.
// @Description get logged-in user data.
// @Tags Users
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=forms.WhoAmIResponse}
// @Failure 500 {object} utils.Response{data=object}
// @Security ApiKeyAuth
// @Router /api/v1/users/me [get]
func (ac *UserController) Me(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorization_payload").(*utils.TokenPayload)

	var user models.User
	err := ac.s.UserService.FindOne(&user, authPayload.UserId, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData(utils.ResponseStatusSuccess, "success get current user", forms.WhoAmIResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: string(user.Status),
	}))
}
