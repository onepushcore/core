package service

import (
	"context"
	"github.com/onepushcore/core"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

const (
	BucketKeyOwnedApps = "onepush::core::ownedapps"
)

type AppConfigService struct {
	bucket *BucketService[core.AppConfig]
}

func NewAppConfigService(rds redis.Cmdable) *AppConfigService {
	return &AppConfigService{
		bucket: &BucketService[core.AppConfig]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return "onepush::core::applications::" + owner
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
	// Config
	account, appKey := config.Account, config.AppKey
	err := s.bucket.Store(ctx, account, appKey, config)
	if err != nil {
		slog.Error("redis store app config error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	// Owned
	ret := s.bucket.Rds.HSet(ctx, BucketKeyOwnedApps, appKey, account)
	if ret.Err() != nil {
		slog.Error("redis store app owner error.", "account", account, "appKey", appKey, "error", ret.Err())
		return ret.Err()
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
	// Owned
	ret := s.bucket.Rds.HExists(ctx, BucketKeyOwnedApps, appKey)
	if ret.Err() != nil {
		slog.Error("redis exists app owner error.", "account", account, "appKey", appKey, "error", ret.Err())
		return false
	}
	if !ret.Val() {
		return false
	}
	// Config
	exists, err := s.bucket.Exists(ctx, account, appKey)
	if err != nil {
		slog.Error("redis exists app config error.", "account", account, "appKey", appKey, "error", err)
		return false
	}
	return exists
}

func (s *AppConfigService) Remove(ctx context.Context, account, appKey string) error {
	// Config
	err := s.bucket.Remove(ctx, account, appKey)
	if err != nil {
		slog.Error("redis remove app config error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	// Owned
	ret := s.bucket.Rds.HDel(ctx, BucketKeyOwnedApps, appKey)
	if ret.Err() != nil {
		slog.Error("redis remove app owner error.", "account", account, "appKey", appKey, "error", ret.Err())
		return ret.Err()
	}
	return nil
}

func (s *AppConfigService) Owner(ctx context.Context, appKey string) (string, error) {
	resp := s.bucket.Rds.HGet(ctx, BucketKeyOwnedApps, appKey)
	if resp.Err() != nil {
		slog.Error("redis get app owner error.", "appKey", appKey, "error", resp.Err())
		return "", resp.Err()
	}
	return resp.Val(), nil
}
