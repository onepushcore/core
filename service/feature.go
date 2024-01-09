package service

import (
	"context"
	"fmt"
	"github.com/onepushcore/core"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type FeatureConfigService struct {
	bucket *BucketService[core.FeatureConfig]
}

func NewFeatureConfigService(rds redis.Cmdable) *FeatureConfigService {
	return &FeatureConfigService{
		bucket: &BucketService[core.FeatureConfig]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return fmt.Sprintf("onepush::core::features::%s", owner)
			},
			NewValue: func() core.FeatureConfig {
				return core.FeatureConfig{}
			},
		},
	}
}

func (s *FeatureConfigService) List(ctx context.Context, appKey string, channelType core.ChannelType) (map[string]core.FeatureConfig, error) {
	dmap, err := s.bucket.List(ctx, s.newOwnerKey(appKey, channelType))
	if err != nil {
		slog.Error("redis list feature config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *FeatureConfigService) Store(ctx context.Context, config core.FeatureConfig) error {
	appKey, channelType := config.AppKey, config.Channel
	err := s.bucket.Store(ctx, s.newOwnerKey(appKey, channelType), string(config.Type), config)
	if err != nil {
		slog.Error("redis store feature config error.", "appKey", appKey, "channelType", channelType, "featureType", config.Type, "error", err)
		return err
	}
	return nil
}

func (s *FeatureConfigService) Load(ctx context.Context, appKey string, channelType core.ChannelType, featureType core.FeatureType) (*core.FeatureConfig, error) {
	value, err := s.bucket.Load(ctx, s.newOwnerKey(appKey, channelType), string(featureType))
	if err != nil {
		slog.Error("redis load feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *FeatureConfigService) Exists(ctx context.Context, appKey string, channelType core.ChannelType, featureType core.FeatureType) bool {
	exists, err := s.bucket.Exists(ctx, s.newOwnerKey(appKey, channelType), string(featureType))
	if err != nil {
		slog.Error("redis exists feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return false
	}
	return exists
}

func (s *FeatureConfigService) Remove(ctx context.Context, appKey string, channelType core.ChannelType, featureType core.FeatureType) error {
	err := s.bucket.Remove(ctx, s.newOwnerKey(appKey, channelType), string(featureType))
	if err != nil {
		slog.Error("redis remove feature config error.", "appKey", appKey, "channelType", channelType, "featureType", featureType, "error", err)
		return err
	}
	return nil
}

func (s *FeatureConfigService) newOwnerKey(appKey string, channelType core.ChannelType) string {
	return appKey + "#" + string(channelType)
}
