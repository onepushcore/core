package service

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/onepushcore/core"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type ChannelConfigService struct {
	rds redis.Cmdable
}

func NewChannelConfigService(rds redis.Cmdable) *ChannelConfigService {
	return &ChannelConfigService{
		rds: rds,
	}
}

func (s *ChannelConfigService) Load(ctx context.Context, channelType core.ChannelType, appKey string) (*core.ChannelConfig, error) {
	bucket := s.newBucket(channelType)
	resp := s.rds.HGet(ctx, bucket, appKey)
	if resp.Err() != nil {
		slog.Error("redis load channel config bucket error.", "bucket", bucket, "appKey", appKey, "error", resp.Err())
		return nil, fmt.Errorf("load channel entity %w", resp.Err())
	}
	if resp.Val() != "" {
		var config = &core.ChannelConfig{}
		err := sonic.ConfigFastest.Unmarshal([]byte(resp.Val()), config)
		if err != nil {
			slog.Error("unmarshal channel config error.", "value", resp.Val(), "error", resp.Err())
			return nil, fmt.Errorf("unmarahsl channel config. %w", err)
		}
		return config, nil
	}
	return nil, nil
}

func (s *ChannelConfigService) Exists(ctx context.Context, channelType core.ChannelType, appKey string) bool {
	bucket := s.newBucket(channelType)
	ret := s.rds.HExists(ctx, bucket, appKey)
	if ret.Err() != nil {
		slog.Error("redis channel config exists error.", "bucket", bucket, "appKey", appKey, "error", ret.Err())
		return false
	}
	return ret.Val()
}

func (s *ChannelConfigService) newBucket(channelType core.ChannelType) string {
	return fmt.Sprintf("onepush::core::channels::%s", channelType)
}
