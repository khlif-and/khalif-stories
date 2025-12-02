package main

import (
	"log"
	// Hapus import meilisearch dan pkg/search dari sini jika tidak dipakai fungsi lain

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/database"
	"khalif-stories/pkg/utils"

)

func ProvideDB(cfg *config.Config) *gorm.DB {
	database.EnsureDBExists(cfg.DBUrl)

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
}

// func ProvideMeili... DIHAPUS

func ProvideAzureUploader(cfg *config.Config) *utils.AzureUploader {
	uploader, err := utils.NewAzureUploader(cfg.AzureConnStr, cfg.AzureContainer)
	if err != nil {
		log.Fatal(err)
	}
	return uploader
}