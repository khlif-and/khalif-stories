package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

)

type RedisRepo struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisRepo) Set(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(context.Background(), key, value, ttl).Err()
}

func (r *RedisRepo) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisRepo) DeletePrefix(prefix string) error {
	ctx := context.Background()
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		r.client.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}