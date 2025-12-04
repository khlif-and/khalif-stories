package mocks

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"

	"khalif-stories/internal/domain"

)

type CategoryUseCaseMock struct {
	mock.Mock
}

func (m *CategoryUseCaseMock) Create(ctx context.Context, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	args := m.Called(ctx, name, file, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) GetAll(ctx context.Context) ([]domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) Get(ctx context.Context, uuid string) (*domain.Category, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) Search(ctx context.Context, query string) ([]domain.Category, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) Update(ctx context.Context, uuid string, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	args := m.Called(ctx, uuid, name, file, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

type StoryUseCaseMock struct {
	mock.Mock
}

func (m *StoryUseCaseMock) Create(ctx context.Context, title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	args := m.Called(ctx, title, desc, categoryID, file, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Story), args.Error(1)
}

func (m *StoryUseCaseMock) GetAll(ctx context.Context, page, limit int, sort string) ([]domain.Story, error) {
	args := m.Called(ctx, page, limit, sort)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (m *StoryUseCaseMock) Search(ctx context.Context, query string) (*[]domain.Story, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]domain.Story), args.Error(1)
}

func (m *StoryUseCaseMock) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *StoryUseCaseMock) AddSlide(ctx context.Context, storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	args := m.Called(ctx, storyUUID, content, sequence, file, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Slide), args.Error(1)
}