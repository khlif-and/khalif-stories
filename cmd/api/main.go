package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/pkg/database"

)

type App struct {
	DB             *gorm.DB
	RDB            *redis.Client
	Meili          *meilisearch.Client
	StoryHandler   *handler.StoryHandler
	ChapterHandler *handler.ChapterHandler
	SearchHandler  *handler.SearchHandler
}

func NewApp(db *gorm.DB, rdb *redis.Client, meili *meilisearch.Client, h *handler.StoryHandler, ch *handler.ChapterHandler, sh *handler.SearchHandler) *App {
	return &App{
		DB:             db,
		RDB:            rdb,
		Meili:          meili,
		StoryHandler:   h,
		ChapterHandler: ch,
		SearchHandler:  sh,
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
		fmt.Println("Resetting Schema...")
		database.ResetSchema(app.DB)
	}

	app.DB.AutoMigrate(&domain.Category{}, &domain.Story{}, &domain.Slide{})

	database.SeedCategories(app.DB)

	r := gin.Default()

	SetupRoutes(r, app, cfg)

	r.Run(":" + cfg.Port)
}