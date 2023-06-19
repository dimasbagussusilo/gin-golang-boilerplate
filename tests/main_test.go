package tests

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/middlewares"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Server serves HTTP requests for our marketplace service.
type Server struct {
	config     config.Config
	store      *gorm.DB
	tokenMaker utils.TokenMaker
	router     *gin.Engine
}

func (server *Server) setupRouter() {
	r := gin.Default()
	authorized := r.Group("/").Use(middlewares.AuthMiddleware(server.tokenMaker))
	authorized.GET("/", func(context *gin.Context) {

	})

	server.router = r
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Config, store *gorm.DB) (*Server, error) {
	tokenMaker, err := utils.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	//if server.config.Environment == 'development' {
	//	gin.SetMode()
	//}
	server.setupRouter()
	return server, nil
}

func newTestServer(t *testing.T, store *gorm.DB) *Server {
	cfg := config.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(cfg, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
