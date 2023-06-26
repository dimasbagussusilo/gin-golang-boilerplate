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
	"strconv"
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
	err := ac.s.UserService.FindOne(ctx, &user, authPayload.UserId, nil)
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

// GetAllUsers godoc
// @Summary Show all user.
// @Description get all users' data.
// @Tags Users
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=object}
// @Failure 500 {object} utils.Response{data=object}
// @Security ApiKeyAuth
// @Param   page     	query     int     false  "Page"     	default(1)
// @Param   limit    	query     int	  false  "Page Limit"   default(10)
// @Param   search   	query     string  false  "Search"
// @Param   email		query     string  false  "Email"
// @Param   order_by	query     string  false  "Order by"
// @Router /api/v1/users [get]
func (ac *UserController) GetAllUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	filterQuery := map[string]any{}
	if email := ctx.Query("email"); email != "" {
		filterQuery["email"] = ctx.Query("email")
	}

	searchQuery := []string{
		"name",
		"email",
	}

	orderQuery := ctx.Query("order_by")

	queryOptions := &service.OptionQuery{
		Page:   page,
		Limit:  limit,
		Filter: filterQuery,
		Search: searchQuery,
		Order:  orderQuery,
	}

	var users []models.User
	pagination, err := ac.s.UserService.FindAll(ctx, &users, queryOptions, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	res := forms.GetAllUserResponse{
		Users:      users,
		Pagination: pagination,
	}

	ctx.JSON(http.StatusOK, utils.ResponseData(utils.ResponseStatusSuccess, "success get current user", res))
}
