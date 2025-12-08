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

		repository.NewCategoryRepository,
		repository.NewStoryRepository,
		repository.NewChapterRepository,
		repository.NewCacheRepository,
		repository.NewPreferenceRepository,

		wire.Bind(new(domain.CategoryRepository), new(*repository.CategoryRepo)),
		wire.Bind(new(domain.StoryRepository), new(*repository.StoryRepo)),
		wire.Bind(new(domain.ChapterRepository), new(*repository.ChapterRepo)),
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepo)),
		wire.Bind(new(domain.PreferenceRepository), new(*repository.PreferenceRepo)),

		usecase.NewCategoryUseCase,
		usecase.NewStoryUseCase,
		usecase.NewChapterUseCase,
		usecase.NewPreferenceUseCase,

		wire.Bind(new(domain.CategoryUseCase), new(*usecase.CategoryUC)),
		wire.Bind(new(domain.ChapterUseCase), new(*usecase.ChapterUC)),
		wire.Bind(new(domain.PreferenceUseCase), new(*usecase.PreferenceUC)),

		handler.NewCategoryHandler,
		handler.NewStoryHandler,
		handler.NewChapterHandler,
		handler.NewPreferenceHandler,

		NewApp,
	)
	return &App{}, nil
}