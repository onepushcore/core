package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type AppEntity struct {
	Account     string         `json:"account"`     // 所属帐号
	AppKey      string         `json:"app_key"`     // 应用唯一标识
	Name        string         `json:"name"`        // 应用名称
	Description string         `json:"description"` // 应用描述
	Attrs       map[string]any `json:"attrs"`       // 应用属性
}

////

const (
	BucketKeyOwnedApps = "onepush::core::ownedapps"
)

type AppService struct {
	bucket *BucketService[AppEntity]
}

func NewAppService(rds redis.Cmdable) *AppService {
	return &AppService{
		bucket: &BucketService[AppEntity]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return "onepush::core::applications::" + owner
			},
			NewValue: func() AppEntity {
				return AppEntity{}
			},
		},
	}
}

func (s *AppService) List(ctx context.Context, account string) (map[string]AppEntity, error) {
	dmap, err := s.bucket.List(ctx, account)
	if err != nil {
		slog.Error("redis list app config error.", "account", account, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *AppService) Store(ctx context.Context, config AppEntity) error {
	// Config
	account, appKey := config.Account, config.AppKey
	err := s.bucket.Store(ctx, account, appKey, config)
	if err != nil {
		slog.Error("redis store app config error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	// Owned
	resp := s.bucket.Rds.HSet(ctx, BucketKeyOwnedApps, appKey, account)
	if err := resp.Err(); err != nil {
		slog.Error("redis store app owner error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	return nil
}

func (s *AppService) Load(ctx context.Context, account, appKey string) (*AppEntity, error) {
	value, err := s.bucket.Load(ctx, account, appKey)
	if err != nil {
		slog.Error("redis load app config error.", "account", account, "appKey", appKey, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *AppService) Exists(ctx context.Context, account, appKey string) bool {
	// Owned
	resp := s.bucket.Rds.HExists(ctx, BucketKeyOwnedApps, appKey)
	if err := resp.Err(); err != nil {
		slog.Error("redis exists app owner error.", "account", account, "appKey", appKey, "error", err)
		return false
	}
	if !resp.Val() {
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

func (s *AppService) Remove(ctx context.Context, account, appKey string) error {
	// Config
	err := s.bucket.Remove(ctx, account, appKey)
	if err != nil {
		slog.Error("redis remove app config error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	// Owned
	resp := s.bucket.Rds.HDel(ctx, BucketKeyOwnedApps, appKey)
	if err := resp.Err(); err != nil {
		slog.Error("redis remove app owner error.", "account", account, "appKey", appKey, "error", err)
		return err
	}
	return nil
}

func (s *AppService) Owner(ctx context.Context, appKey string) (string, error) {
	resp := s.bucket.Rds.HGet(ctx, BucketKeyOwnedApps, appKey)
	if err := resp.Err(); err != nil {
		slog.Error("redis get app owner error.", "appKey", appKey, "error", err)
		return "", err
	}
	return resp.Val(), nil
}
