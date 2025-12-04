package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/pkg/database"

)

type App struct {
	DB              *gorm.DB
	RDB             *redis.Client
	CategoryHandler *handler.CategoryHandler
	StoryHandler    *handler.StoryHandler
}

func NewApp(db *gorm.DB, rdb *redis.Client, ch *handler.CategoryHandler, sh *handler.StoryHandler) *App {
	return &App{
		DB:              db,
		RDB:             rdb,
		CategoryHandler: ch,
		StoryHandler:    sh,
	}
}

func main() {
	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()
	app, err := InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	if *refreshFlag {
		database.ResetSchema(app.DB)
	}

	app.DB.AutoMigrate(&domain.Category{}, &domain.Story{}, &domain.Slide{})
	
	database.SetupDatabaseCapabilities(app.DB)

	database.SeedCategories(app.DB)

	r := gin.Default()
	SetupRoutes(r, app, cfg)
	r.Run(":" + cfg.Port)
}