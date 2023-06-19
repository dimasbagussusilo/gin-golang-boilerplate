package controllers

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServerController struct {
	config *config.Config
	db     *gorm.DB
}

func NewServerController(config *config.Config, db *gorm.DB) *ServerController {
	return &ServerController{
		config: config,
		db:     db,
	}
}

func (s *ServerController) HealthCheck(c *gin.Context) {
	// Server Check
	dbConn, _ := s.db.DB()

	parsedDsn, _ := url.Parse(s.config.DBSource)
	dbHost := parsedDsn.Host
	dbName := parsedDsn.Path

	if dbHost == "" {
		// Parse DSN server format
		pairs := strings.Split(dbName, " ")
		serverData := make(map[string]string)
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				serverData[parts[0]] = parts[1]
			}
		}
		dbHost = serverData["host"] + ":" + serverData["port"]
		dbName = serverData["dbname"]
	}

	var databaseStatus string
	if err := dbConn.Ping(); err != nil {
		databaseStatus = "error"
	} else {
		databaseStatus = "ok"
	}

	c.JSON(http.StatusOK, utils.ResponseData(utils.ResponseStatusSuccess, "Server running well", map[string]any{
		"server_status":   "ok",
		"database_status": databaseStatus,
		"database_name":   dbName,
		"database_host":   dbHost,
	},
	))
}
