package usertoken

import (
	"context"
	"fmt"
	"time"

	"code-search/auth-service/entity/config"
	"github.com/go-redis/redis/v8"
)

func NewClient() Client {
	cfg := config.DefaultConfig
	// 连接Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return &impl{
		redis: rdb,
	}
}

type Client interface {
	SetUserToken(ctx context.Context, token, username string) error
	ValidateToken(ctx context.Context, token string) (string, error)
}

type impl struct {
	redis *redis.Client
}

func (i *impl) SetUserToken(ctx context.Context, token, username string) error {
	err := i.redis.Set(ctx, token, username, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("set user %s token %s error", username, token)
	}
	return nil
}

func (i *impl) ValidateToken(ctx context.Context, token string) (string, error) {
	res, err := i.redis.Get(ctx, token).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
