package app

import (
	"gosearch/internal/config"
	"gosearch/internal/store"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Config *config.Config
}

func InitContainer() *Container {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := store.InitDB(cfg.Database.DSN)
	if err != nil {
		panic(err)
	}

	return &Container{
		DB:     db,
		Config: cfg,
	}
}
