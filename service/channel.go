package service

import (
	"context"
	"github.com/onepushcore/core"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type ChannelConfigService struct {
	bucket *BucketService[core.ChannelConfig]
}

func NewChannelConfigService(rds redis.Cmdable) *ChannelConfigService {
	return &ChannelConfigService{
		bucket: &BucketService[core.ChannelConfig]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return "onepush::core::channels::" + owner
			},
			NewValue: func() core.ChannelConfig {
				return core.ChannelConfig{}
			},
		},
	}
}

func (s *ChannelConfigService) List(ctx context.Context, appKey string) (map[string]core.ChannelConfig, error) {
	dmap, err := s.bucket.List(ctx, appKey)
	if err != nil {
		slog.Error("redis list channel config error.", "appKey", appKey, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *ChannelConfigService) Store(ctx context.Context, config core.ChannelConfig) error {
	err := s.bucket.Store(ctx, config.AppKey, string(config.Type), config)
	if err != nil {
		slog.Error("redis store channel config error.", "appKey", config.AppKey, "channelType", config.Type, "error", err)
		return err
	}
	return nil
}

func (s *ChannelConfigService) Load(ctx context.Context, appKey string, channelType core.ChannelType) (*core.ChannelConfig, error) {
	value, err := s.bucket.Load(ctx, appKey, string(channelType))
	if err != nil {
		slog.Error("redis load channel config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *ChannelConfigService) Exists(ctx context.Context, appKey string, channelType core.ChannelType) bool {
	exists, err := s.bucket.Exists(ctx, appKey, string(channelType))
	if err != nil {
		slog.Error("redis exists channel config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return false
	}
	return exists
}

func (s *ChannelConfigService) Remove(ctx context.Context, appKey string, channelType core.ChannelType) error {
	err := s.bucket.Remove(ctx, appKey, string(channelType))
	if err != nil {
		slog.Error("redis remove channel config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return err
	}
	return nil
}
