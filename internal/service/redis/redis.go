package redis

import (
	"context"
	"time"

	"shopnexus-go-service/config"
	redisclient "shopnexus-go-service/internal/client/redis"
)

type Service interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type ServiceImpl struct {
	client redisclient.Client
}

func NewService() (Service, error) {
	client, err := redisclient.NewClient(redisclient.RedisConfig{
		Addr:     []string{config.GetConfig().Redis.Addr},
		Password: config.GetConfig().Redis.Password,
		DB:       config.GetConfig().Redis.DB,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{
		client: client,
	}, nil
}

func (r *ServiceImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *ServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

func (r *ServiceImpl) Delete(ctx context.Context, key string) error {
	return r.client.Delete(ctx, key)
}

func (r *ServiceImpl) Exists(ctx context.Context, key string) (bool, error) {
	return r.client.Exists(ctx, key)
}
