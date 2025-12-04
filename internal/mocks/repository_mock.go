package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"khalif-stories/internal/domain"

)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetByUUID(ctx context.Context, uuid string) (*domain.Category, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetAll(ctx context.Context) ([]domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) Search(ctx context.Context, query string) ([]domain.Category, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) Update(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) UpdateColor(ctx context.Context, id uint, color string) error {
	args := m.Called(ctx, id, color)
	return args.Error(0)
}

type StoryRepositoryMock struct {
	mock.Mock
}

func (m *StoryRepositoryMock) Create(ctx context.Context, story *domain.Story) error {
	args := m.Called(ctx, story)
	return args.Error(0)
}

func (m *StoryRepositoryMock) GetByID(ctx context.Context, id uint) (*domain.Story, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Story), args.Error(1)
}

func (m *StoryRepositoryMock) GetByUUID(ctx context.Context, uuid string) (*domain.Story, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Story), args.Error(1)
}

func (m *StoryRepositoryMock) GetAll(ctx context.Context, page, limit int, sort string) ([]domain.Story, error) {
	args := m.Called(ctx, page, limit, sort)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (m *StoryRepositoryMock) Search(ctx context.Context, query string) ([]domain.Story, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (m *StoryRepositoryMock) Update(ctx context.Context, story *domain.Story) error {
	args := m.Called(ctx, story)
	return args.Error(0)
}

func (m *StoryRepositoryMock) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *StoryRepositoryMock) UpdateColor(ctx context.Context, id uint, color string) error {
	args := m.Called(ctx, id, color)
	return args.Error(0)
}

func (m *StoryRepositoryMock) CreateSlide(ctx context.Context, slide *domain.Slide) error {
	args := m.Called(ctx, slide)
	return args.Error(0)
}

func (m *StoryRepositoryMock) CountSlides(ctx context.Context, storyID uint) (int64, error) {
	args := m.Called(ctx, storyID)
	return args.Get(0).(int64), args.Error(1)
}

type RedisRepositoryMock struct {
	mock.Mock
}

func (m *RedisRepositoryMock) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *RedisRepositoryMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *RedisRepositoryMock) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisRepositoryMock) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}