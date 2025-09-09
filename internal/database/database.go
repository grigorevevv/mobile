package database

import (
	"fmt"
	"log"
	config "mobile/internal/configs"
	"mobile/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBase struct {
	DB *gorm.DB
}

func NewDataBase(c *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connection Database")
		return nil, err
	}

	if err := db.AutoMigrate(
		domain.Subscription{},
	); err != nil {
		return db, err
	}

	log.Printf("Database connection was created: %s \n", c.DBName)
	return db, nil
}
