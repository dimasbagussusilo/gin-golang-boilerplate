package seeders

import (
	"encoding/json"
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
	"os"
	"time"
)

func (s Seed) User(c *config.Config, srv *service.Services) {
	var jsonFile []byte
	var err error

	if c.Environment == "development" {
		jsonFile, err = os.ReadFile("database/seeders/seeds/user.json")
	} else {
		jsonFile, err = os.ReadFile("database/seeders/seeds/user_prod.json")
	}
	if err != nil {
		fmt.Println("error when parse json file", err)
	}

	var users []models.User
	err = json.Unmarshal(jsonFile, &users)
	if err != nil {
		fmt.Println("error while parse json file")
		return
	}

	var totalUser int64
	srv.UserService.CustomQuery().Count(&totalUser)

	if totalUser == 0 {
		for _, user := range users {
			userParams := &models.User{
				Name:      user.Name,
				Email:     user.Email,
				Password:  user.Password,
				Status:    user.Status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, err = srv.UserService.Create(userParams, nil)
			if err != nil {
				fmt.Println("error while create seed data")
				return
			}
		}
		srv.UserService.CustomQuery().Count(&totalUser)
		if err != nil {
			fmt.Println("error when get totalUser", err)
		} else if int(totalUser) < len(users) {
			fmt.Println("User seeding incomplete")
		} else {
			fmt.Println("User seeding successful")
		}
	}

}
