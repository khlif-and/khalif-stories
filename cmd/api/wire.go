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
	// HAPUS BARIS INI: "khalif-stories/pkg/utils" 

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
		repository.NewChapterRepository,
		repository.NewCacheRepository,

		// Binding Repository Interface -> Implementation
		wire.Bind(new(domain.CategoryRepository), new(*repository.CategoryRepo)),
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),
		wire.Bind(new(domain.ChapterRepository), new(*repository.ChapterRepo)),
		
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepo)),

		// UseCase (Provider)
		usecase.NewCategoryUseCase,
		usecase.NewStoryUseCase,
		usecase.NewChapterUseCase,

		// Binding Interface UseCase -> Implementation Struct
		wire.Bind(new(domain.CategoryUseCase), new(*usecase.CategoryUC)),
		wire.Bind(new(domain.ChapterUseCase), new(*usecase.ChapterUC)),

		// Handler
		handler.NewCategoryHandler,
		handler.NewStoryHandler,
		handler.NewChapterHandler,

		NewApp,
	)
	return &App{}, nil
}