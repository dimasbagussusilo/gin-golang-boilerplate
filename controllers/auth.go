package controllers

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
// @Success 200 {object} utils.Response{data=forms.SignupResponse}
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

		createdUser, err := ac.s.UserService.Create(ctx, &models.User{
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

	res := forms.SignupResponse{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	ctx.JSON(http.StatusCreated, utils.ResponseData(utils.ResponseStatusSuccess, "success create user", res))
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
func (ac *AuthController) Signin(ctx *gin.Context) {
	cfg, _ := config.LoadConfig(".")
	tokenMaker, _ := utils.NewJWTMaker(cfg.TokenSymmetricKey)

	var input forms.SigninRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		return
	}

	filterQuery := map[string]any{
		"custom_fields": []clause.Expr{
			gorm.Expr("LOWER(email) LIKE ?", strings.ToLower(input.Email)),
		},
	}

	optionsQuery := &service.OptionQuery{
		Filter: filterQuery,
		Search: nil,
		Page:   1,
		Limit:  1,
	}

	var users []models.User
	_, err := ac.s.UserService.FindAll(ctx, &users, optionsQuery, nil)
	if err != nil {
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "invalid email or password", nil))
		return
	}

	var user models.User
	if len(users) > 0 {
		user = users[0]
	}

	//if user.Status == models.UserStatusPending {
	//	ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "please verify your account", nil))
	//	return
	//}

	err = utils.CheckPassword(input.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData(utils.ResponseStatusError, "invalid email or password", nil))
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

	_, err = ac.s.TokenService.Create(ctx, &models.Token{
		ID:           refreshPayload.Id,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    time.Time{},
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
	}

	rsp := forms.SigninResponse{
		SessionID:             refreshPayload.Id,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusCreated, utils.ResponseData(utils.ResponseStatusSuccess, "success signin user", rsp))
}
