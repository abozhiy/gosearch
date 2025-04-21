package store

import (
	"fmt"

	"gosearch/internal/store/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DSN  string `yaml:"dsn"`
	Pool int    `yaml:"pool"`
}

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB > %v", err)
	}

	if err := db.AutoMigrate(migrations.Models()...); err != nil {
		return nil, fmt.Errorf("migration failed > %w", err)
	}

	return db, nil
}
