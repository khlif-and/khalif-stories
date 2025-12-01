//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/internal/repository"
	"khalif-stories/internal/usecase"

)

type App struct {
	DB             *gorm.DB
	RDB            *redis.Client
	StoryHandler   *handler.StoryHandler
	ChapterHandler *handler.ChapterHandler
}

func NewApp(db *gorm.DB, rdb *redis.Client, h *handler.StoryHandler, ch *handler.ChapterHandler) *App {
	return &App{
		DB:             db,
		RDB:            rdb,
		StoryHandler:   h,
		ChapterHandler: ch,
	}
}

func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,
		ProvideAzureUploader,

		repository.NewStoryRepository,
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),

		usecase.NewStoryUseCase,

		handler.NewStoryHandler,
		handler.NewChapterHandler,

		NewApp,
	)
	return &App{}, nil
}