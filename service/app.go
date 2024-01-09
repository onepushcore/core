package service

import (
	"context"
	"fmt"
	"github.com/onepushcore/core"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type AppConfigService struct {
	bucket *BucketService[core.AppConfig]
}

func NewAppConfigService(rds redis.Cmdable) *AppConfigService {
	return &AppConfigService{
		bucket: &BucketService[core.AppConfig]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return fmt.Sprintf("onepush::core::applications::%s", owner)
			},
			NewValue: func() core.AppConfig {
				return core.AppConfig{}
			},
		},
	}
}

func (s *AppConfigService) List(ctx context.Context, account string) (map[string]core.AppConfig, error) {
	dmap, err := s.bucket.List(ctx, account)
	if err != nil {
		slog.Error("redis list app config error.", "account", account, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *AppConfigService) Store(ctx context.Context, config core.AppConfig) error {
	account, appKey := config.Account, config.AppKey
	err := s.bucket.Store(ctx, account, appKey, config)
	if err != nil {
		slog.Error("redis store app config error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	return nil
}

func (s *AppConfigService) Load(ctx context.Context, account, appKey string) (*core.AppConfig, error) {
	value, err := s.bucket.Load(ctx, account, appKey)
	if err != nil {
		slog.Error("redis load app config error.", "account", account, "appKey", appKey, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *AppConfigService) Exists(ctx context.Context, account, appKey string) bool {
	exists, err := s.bucket.Exists(ctx, account, appKey)
	if err != nil {
		slog.Error("redis exists app config error.", "account", account, "appKey", appKey, "error", err)
		return false
	}
	return exists
}
