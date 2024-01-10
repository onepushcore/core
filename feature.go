package core

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type FeatureEntity struct {
	AppKey      string         `json:"app_key"`      // 所属应用标识
	ChannelType string         `json:"channel_type"` // 所属渠道类型
	Type        string         `json:"type"`         // 特性类型
	Enabled     bool           `json:"enabled"`      // 是否启用
	Attrs       map[string]any `json:"attrs"`        // 属性列表
	Options     map[string]any `json:"options"`      // 配置参数
}

////

type FeatureService struct {
	bucket *BucketService[FeatureEntity]
}

func NewFeatureService(rds redis.Cmdable) *FeatureService {
	return &FeatureService{
		bucket: &BucketService[FeatureEntity]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return fmt.Sprintf("onepush::core::features::%s", owner)
			},
			NewValue: func() FeatureEntity {
				return FeatureEntity{}
			},
		},
	}
}

func (s *FeatureService) List(ctx context.Context, appKey string, channelType string) (map[string]FeatureEntity, error) {
	dmap, err := s.bucket.List(ctx, s.newOwnerKey(appKey, channelType))
	if err != nil {
		slog.Error("redis list feature config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *FeatureService) Store(ctx context.Context, feature FeatureEntity) error {
	appKey, channelType := feature.AppKey, feature.ChannelType
	err := s.bucket.Store(ctx, s.newOwnerKey(appKey, channelType), feature.Type, feature)
	if err != nil {
		slog.Error("redis store feature feature error.", "appKey", appKey, "channelType", channelType, "featureType", feature.Type, "error", err)
		return err
	}
	return nil
}

func (s *FeatureService) Load(ctx context.Context, appKey string, channelType string, featureType string) (*FeatureEntity, error) {
	value, err := s.bucket.Load(ctx, s.newOwnerKey(appKey, channelType), featureType)
	if err != nil {
		slog.Error("redis load feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *FeatureService) Exists(ctx context.Context, appKey string, channelType string, featureType string) bool {
	exists, err := s.bucket.Exists(ctx, s.newOwnerKey(appKey, channelType), featureType)
	if err != nil {
		slog.Error("redis exists feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return false
	}
	return exists
}

func (s *FeatureService) Remove(ctx context.Context, appKey string, channelType string, featureType string) error {
	err := s.bucket.Remove(ctx, s.newOwnerKey(appKey, channelType), featureType)
	if err != nil {
		slog.Error("redis remove feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return err
	}
	return nil
}

func (s *FeatureService) newOwnerKey(appKey string, channelType string) string {
	return appKey + "#" + channelType
}
