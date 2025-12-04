package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/mocks"
	"khalif-stories/internal/usecase"

)

func TestStoryUseCase_Create(t *testing.T) {
	mockRepo := new(mocks.StoryRepositoryMock)
	cfg := &config.Config{SlideLimit: 20}

	uc := usecase.NewStoryUseCase(cfg, mockRepo, nil, nil)

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Story")).Return(nil)

		res, err := uc.Create(ctx, "Title", "Desc", 1, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Title", res.Title)
		assert.Equal(t, "Draft", res.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestStoryUseCase_AddSlide(t *testing.T) {
	mockRepo := new(mocks.StoryRepositoryMock)
	cfg := &config.Config{SlideLimit: 5}

	uc := usecase.NewStoryUseCase(cfg, mockRepo, nil, nil)
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		storyUUID := "abc-123"
		story := &domain.Story{ID: 1, UUID: storyUUID}

		mockRepo.On("GetByUUID", ctx, storyUUID).Return(story, nil)
		mockRepo.On("CountSlides", ctx, uint(1)).Return(int64(2), nil)
		mockRepo.On("CreateSlide", ctx, mock.AnythingOfType("*domain.Slide")).Return(nil)

		res, err := uc.AddSlide(ctx, storyUUID, "Content", 1, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("limit reached", func(t *testing.T) {
		storyUUID := "abc-999"
		story := &domain.Story{ID: 2, UUID: storyUUID}

		mockRepo.On("GetByUUID", ctx, storyUUID).Return(story, nil)
		mockRepo.On("CountSlides", ctx, uint(2)).Return(int64(5), nil)

		res, err := uc.AddSlide(ctx, storyUUID, "Content", 1, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "slide limit reached", err.Error())
	})
}