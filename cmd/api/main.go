package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/pkg/database"
	"khalif-stories/pkg/logger"

)

// @title           Khalif Stories API
// @version         1.0
// @description     API Service for Khalif Stories Application
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@khalifstories.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type App struct {
	DB              *gorm.DB
	RDB             *redis.Client
	CategoryHandler *handler.CategoryHandler
	StoryHandler    *handler.StoryHandler
	ChapterHandler  *handler.ChapterHandler
}

func NewApp(db *gorm.DB, rdb *redis.Client, ch *handler.CategoryHandler, sh *handler.StoryHandler, chapH *handler.ChapterHandler) *App {
	return &App{
		DB:              db,
		RDB:             rdb,
		CategoryHandler: ch,
		StoryHandler:    sh,
		ChapterHandler:  chapH,
	}
}

func main() {
	logger.Init()

	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()
	app, err := InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}

	if *refreshFlag {
		database.ResetSchema(app.DB)
		logger.Info("Database reset successfully")
	}

	app.DB.AutoMigrate(&domain.Category{}, &domain.Story{}, &domain.Chapter{}, &domain.Slide{})

	database.RunMigrations(app.DB)

	database.SeedCategories(app.DB)

	r := gin.New()
	r.Use(gin.Recovery())

	SetupRoutes(r, app, cfg)

	logger.Info("Server starting", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Server start failed", zap.Error(err))
	}
}