package routers

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/controllers"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/middlewares"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// V1 Routes Docs
// @title Filesystem API
// @version 1.0
// @description This server provides the Filesystem API needs.
// @termsOfService http://swagger.io/terms/
// @contact.name Dimas Bagus Susilo
// @contact.url http://www.linkedin.com/in/dimasbagussusilo
// @contact.email dimasbagussusilo@gmail.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func V1(r *gin.Engine, c *config.Config, db *gorm.DB, s *service.Services) *gin.Engine {
	// Health Check
	serverController := controllers.NewServerController(c, db)
	r.GET("/health", serverController.HealthCheck)

	//////////
	// Public V1
	v1 := r.Group("api/v1")

	// Auth
	authEndpoint := "/auth"
	auth := controllers.NewAuthController(c, db, s)
	v1.POST(authEndpoint+"/signin", auth.Signin)
	v1.POST(authEndpoint+"/signup", auth.Signup)
	//v1.POST(authEndpoint+"/verify", auth.Verify)

	// User
	usersEndpoint := "/users"
	users := controllers.NewUserController(c, db, s)

	//////////////
	// Authorized V1
	tokenMaker, _ := utils.NewJWTMaker(c.TokenSymmetricKey)
	authorizedV1 := r.Group("api/v1").Use(middlewares.AuthMiddleware(tokenMaker))

	// User
	authorizedV1.GET(usersEndpoint+"/", users.GetAllUsers)
	authorizedV1.GET(usersEndpoint+"/me", users.Me)
	return r
}
