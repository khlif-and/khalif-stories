package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

)

type RedisRepoMock struct {
	mock.Mock
}

func (m *RedisRepoMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *RedisRepoMock) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *RedisRepoMock) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisRepoMock) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}