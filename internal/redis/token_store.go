package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore interface {
	Save(ctx context.Context, customerID, sessionID, token string, ttl time.Duration) error
	Get(ctx context.Context, customerID, sessionID string) (string, error)
	Delete(ctx context.Context, customerID, sessionID string) error
}

type redisTokenStore struct {
	client *redis.Client
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &redisTokenStore{client}
}

func (s *redisTokenStore) Save(ctx context.Context, customerID, sessionID, token string, ttl time.Duration) error {
	return s.client.Set(ctx, refreshKey(customerID, sessionID), token, ttl).Err()
}

func (s *redisTokenStore) Get(ctx context.Context, customerID, sessionID string) (string, error) {
	val, err := s.client.Get(ctx, refreshKey(customerID, sessionID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *redisTokenStore) Delete(ctx context.Context, customerID, sessionID string) error {
	return s.client.Del(ctx, refreshKey(customerID, sessionID)).Err()
}

func refreshKey(customerID, sessionID string) string {
	return "refresh:" + customerID + ":" + sessionID
}
