package main

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/database"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/database/seeders"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/server"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
)

func main() {
	// Load config
	c, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config! %v", err))
	}

	// Initialize database and auto migration
	db, err := database.Init(c)
	if err != nil {
		panic(err)
	}

	// Initialize service
	s := service.Init(db)

	// Initialize seeders
	seeders.Execute(c, s)

	//	Initialize server
	err = server.Init(c, db, s)
	if err != nil {
		panic(err)
	}
}
