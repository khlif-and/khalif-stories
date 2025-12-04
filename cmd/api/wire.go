//go:build wireinject
// +build wireinject

package main

import (
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
		ProvideAzureUploader,

		// Repository (Perlu Bind karena return-nya struct pointer)
		repository.NewCategoryRepository,
		repository.NewStoryRepository,
		repository.NewCacheRepository,

		wire.Bind(new(domain.CategoryRepository), new(*repository.CategoryRepo)),
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),

		// UseCase (TIDAK PERLU Bind karena return-nya sudah interface)
		usecase.NewCategoryUseCase,
		usecase.NewStoryUseCase,

		// Handler
		handler.NewCategoryHandler,
		handler.NewStoryHandler,

		NewApp,
	)
	return &App{}, nil
}