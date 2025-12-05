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
	"khalif-stories/pkg/utils" // [Perlu import ini untuk bind AzureUploader]

)

func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,
		ProvideAzureUploader,

		// Repository (Provider)
		repository.NewCategoryRepository,
		repository.NewStoryRepository,
		repository.NewCacheRepository,

		// Binding Repository Interface -> Implementation
		wire.Bind(new(domain.CategoryRepository), new(*repository.CategoryRepo)),
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),
		
		// [BARU] Binding untuk Dependency CategoryUseCase
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepo)),
		wire.Bind(new(domain.StorageRepository), new(*utils.AzureUploader)),

		// UseCase (Provider)
		usecase.NewCategoryUseCase,
		usecase.NewStoryUseCase,

		// [BARU] Binding Interface UseCase -> Implementation Struct
		// Diperlukan karena NewCategoryUseCase mengembalikan *CategoryUC, tapi Handler minta domain.CategoryUseCase
		wire.Bind(new(domain.CategoryUseCase), new(*usecase.CategoryUC)),

		// Handler
		handler.NewCategoryHandler,
		handler.NewStoryHandler,

		NewApp,
	)
	return &App{}, nil
}