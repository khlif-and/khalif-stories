package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"khalif-stories/internal/domain"
	"khalif-stories/internal/mocks"
	"khalif-stories/internal/usecase"

)

func TestCategoryUseCase_Create(t *testing.T) {
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.CategoryRepositoryMock)
		mockRedis := new(mocks.RedisRepositoryMock)
		
		uc := usecase.NewCategoryUseCase(mockRepo, mockRedis, nil)

		mockRepo.On("GetByName", ctx, "New Category").Return(nil, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Category")).Return(nil)
		
		mockRedis.On("DeletePrefix", ctx, mock.Anything).Return(nil).Maybe()
		mockRedis.On("Del", ctx, mock.Anything).Return(nil).Maybe()

		res, err := uc.Create(ctx, "New Category", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "New Category", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("duplicate name", func(t *testing.T) {
		mockRepo := new(mocks.CategoryRepositoryMock)
		mockRedis := new(mocks.RedisRepositoryMock)
		
		uc := usecase.NewCategoryUseCase(mockRepo, mockRedis, nil)

		existingCategory := &domain.Category{Name: "Existing"}
		mockRepo.On("GetByName", ctx, "Existing").Return(existingCategory, nil)

		res, err := uc.Create(ctx, "Existing", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "category name exists", err.Error())
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.CategoryRepositoryMock)
		mockRedis := new(mocks.RedisRepositoryMock)
		
		uc := usecase.NewCategoryUseCase(mockRepo, mockRedis, nil)

		mockRepo.On("GetByName", ctx, "Error Cat").Return(nil, nil)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		res, err := uc.Create(ctx, "Error Cat", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCategoryUseCase_GetAll(t *testing.T) {
	ctx := context.TODO()

	t.Run("success from db", func(t *testing.T) {
		mockRepo := new(mocks.CategoryRepositoryMock)
		mockRedis := new(mocks.RedisRepositoryMock)
		
		uc := usecase.NewCategoryUseCase(mockRepo, mockRedis, nil)

		categories := []domain.Category{
			{Name: "Cat 1"},
			{Name: "Cat 2"},
		}

		mockRedis.On("Get", ctx, mock.Anything).Return("", errors.New("redis nil")).Maybe()
		mockRedis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
		mockRepo.On("GetAll", ctx).Return(categories, nil)

		res, err := uc.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		mockRepo.AssertExpectations(t)
	})
}