package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore interface {
	Save(ctx context.Context, customerID string, token string, ttl time.Duration) error
	Get(ctx context.Context, customerID string) (string, error)
	Delete(ctx context.Context, customerID string) error
}

type redisTokenStore struct {
	client *redis.Client
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &redisTokenStore{client}
}

func (s *redisTokenStore) Save(ctx context.Context, customerID string, token string, ttl time.Duration) error {
	return s.client.Set(ctx, tokenKey(customerID), token, ttl).Err()
}

func (s *redisTokenStore) Get(ctx context.Context, customerID string) (string, error) {
	val, err := s.client.Get(ctx, tokenKey(customerID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *redisTokenStore) Delete(ctx context.Context, customerID string) error {
	return s.client.Del(ctx, tokenKey(customerID)).Err()
}

func tokenKey(customerID string) string {
	return "token:" + customerID
}
