package database

import (
	"fmt"
	"log"
	"product_service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(cfg config.Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)

	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)

	}

	db, err := gormDb.DB()
	db.SetConnMaxIdleTime(100)

	return gormDb
}
