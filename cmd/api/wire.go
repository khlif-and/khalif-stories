//go:build wireinject
// +build wireinject

package main

import (
	// Hapus import meilisearch disini

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/internal/repository"
	"khalif-stories/internal/usecase"

)

// Kita perlu mendefinisikan struct App lagi disini atau import dari main (karena package sama 'main', aman)
// Tapi agar wire bisa generate, struct App harus terlihat.

func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,
		NewMeiliClientFromConfig, // Panggil wrapper lokal tadi
		ProvideAzureUploader,

		repository.NewStoryRepository,
		repository.NewCacheRepository,
		repository.NewSearchRepository,
		
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),
		wire.Bind(new(domain.SearchRepository), new(*repository.SearchRepo)),

		usecase.NewStoryUseCase,

		handler.NewStoryHandler,
		handler.NewChapterHandler,
		handler.NewSearchHandler,

		NewApp,
	)
	return &App{}, nil
}