package seeders

import (
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/service"
)

type Seed struct {
}

func Execute(c *config.Config, s *service.Services) {
	var seed Seed
	seed.User(c, s)
}
