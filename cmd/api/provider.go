package main

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/database"
	"khalif-stories/pkg/utils"

)

func ProvideDB(cfg *config.Config) *gorm.DB {
	database.EnsureDBExists(cfg.DBUrl)

	dbLogger := logger.Default.LogMode(logger.Error)

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
}

func ProvideAzureUploader(cfg *config.Config) *utils.AzureUploader {
	uploader, err := utils.NewAzureUploader(cfg.AzureConnStr, cfg.AzureContainer)
	if err != nil {
		log.Fatal(err)
	}
	return uploader
}