package controllers

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/database"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/forms"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	c  *config.Config
	db *gorm.DB
	s  *service.Services
}

func NewAuthController(config *config.Config, db *gorm.DB, s *service.Services) *AuthController {
	return &AuthController{
		c:  config,
		db: db,
		s:  s,
	}
}

// Signup godoc
// @Summary Signup user.
// @Description register user.
// @Tags Auth
// @Accept application/json
// @Param request body forms.SignupRequest true "request body"
// @Produce json
// @Success 200 {object} utils.Response{data=object}
// @Failure 400 {object} utils.Response{data=object}
// @Failure 500 {object} utils.Response{data=object}
// @Router /api/v1/auth/signup [post]
func (ac *AuthController) Signup(ctx *gin.Context) {
	var input forms.SignupRequest
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	createUserTransaction := func(tx *gorm.DB) error {
		hashedPassword, _ := utils.HashPassword(input.Password)
		input.Password = hashedPassword

		createdUser, err := ac.s.UserService.Create(&models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
			Status:   models.UserStatusPending,
		}, tx)
		fmt.Println("createdUser", createdUser)
		if err != nil {
			return err
		}

		return nil
	}

	if err := utils.Transaction(database.GetDB(), createUserTransaction); err != nil {
		if err == gorm.ErrDuplicatedKey {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "email already use", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, utils.ResponseData(utils.ResponseStatusSuccess, "success create user", nil))
}

// Signin godoc
// @Summary Login user.
// @Description login user with credentials.
// @Tags Auth
// @Accept application/json
// @Param request body forms.SigninRequest true "request body"
// @Produce json
// @Success 200 {object} utils.Response{data=forms.SigninResponse}
// @Failure 400 {object} utils.Response{data=object}
// @Failure 500 {object} utils.Response{data=object}
// @Router /api/v1/auth/signin [post]
func (ac *AuthController) Signin(c *gin.Context) {
	cfg, _ := config.LoadConfig(".")
	tokenMaker, _ := utils.NewJWTMaker(cfg.TokenSymmetricKey)

	var input forms.SigninRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	queryOptions := func(query *gorm.DB) *gorm.DB {
		// Filter
		query.Where("LOWER(email) = ?", strings.ToLower(input.Email))

		// Sorter
		queryOrder := c.Query("order_by")
		query = utils.SortBy(queryOrder, query)

		return query
	}

	var users []models.User
	_, err := ac.s.UserService.FindAll(&users, 1, 1, queryOptions, nil)
	if err != nil {
		return
	}

	user := users[0]

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "invalid email or password", nil))
		return
	}

	//if user.Status == models.UserStatusPending {
	//	c.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "please verify your account", nil))
	//	return
	//}

	err = utils.CheckPassword(input.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "invalid email or password", nil))
		return
	}

	accessToken, accessPayload, _ := tokenMaker.CreateToken(
		user.ID,
		cfg.AccessTokenDuration,
	)

	refreshToken, refreshPayload, err := tokenMaker.CreateToken(
		user.ID,
		cfg.RefreshTokenDuration,
	)

	_, err = ac.s.TokenService.Create(&models.Token{
		ID:           refreshPayload.Id,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    time.Time{},
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
	}

	rsp := forms.SigninResponse{
		SessionID:             refreshPayload.Id,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	c.JSON(http.StatusCreated, utils.ResponseData(utils.ResponseStatusSuccess, "success signin user", rsp))
}
