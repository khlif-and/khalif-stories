//go:build wireinject
// +build wireinject

package main

import (
	// Hapus baris redis dan gorm di sini

	"github.com/google/wire"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/internal/repository"
	"khalif-stories/internal/usecase"

)

func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,
		NewMeiliClientFromConfig,
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