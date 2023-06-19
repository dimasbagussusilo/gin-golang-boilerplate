package database

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/config"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

var DB *gorm.DB

func Init(c *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DBSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %+e", err)
	}

	err = migrationInit(db)
	if err != nil {
		return nil, fmt.Errorf("error while running auto migrations: %+e", err)
	}

	DB = db

	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}

func migrationInit(db *gorm.DB) error {
	for _, m := range getModels() {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("AutoMigrate failed for model %v: %v\n", reflect.TypeOf(m), err)
		}
	}
	return nil
}

func getModels() []any {
	return []any{
		&models.User{},
		&models.Token{},
	}
}
