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

func (s *ChannelConfigService) List(ctx context.Context, appKey string) (map[core.ChannelType]core.ChannelConfig, error) {
	bucket := s.newBucket(appKey)
	resp := s.rds.HGetAll(ctx, bucket)
	if err := resp.Err(); err != nil {
		slog.Error("redis range channel config error.", "appKey", appKey, "error", err)
		return nil, fmt.Errorf("range channel config. %w", err)
	}
	output := make(map[core.ChannelType]core.ChannelConfig, len(resp.Val()))
	for typ, val := range resp.Val() {
		var config core.ChannelConfig
		if err := sonic.Unmarshal([]byte(val), &config); err != nil {
			slog.Error("unmarshal channel config error.", "appKey", appKey, "config", val, "error", err)
			return nil, fmt.Errorf("unmarshal channel config. %w", err)
		}
		output[core.ChannelType(typ)] = config
	}
	return output, nil
}

func (s *ChannelConfigService) Store(ctx context.Context, config core.ChannelConfig) error {
	appKey, channelType := config.AppKey, config.Type
	data, err := sonic.ConfigFastest.Marshal(config)
	if err != nil {
		slog.Error("marshal channel config error.", "appKey", appKey, "channelType", channelType, "config", config, "error", err)
		return fmt.Errorf("marshal channel config. %w", err)
	}
	bucket := s.newBucket(appKey)
	resp := s.rds.HSet(ctx, bucket, string(channelType), string(data))
	if err := resp.Err(); err != nil {
		slog.Error("redis store channel config error.", "appKey", appKey, "channelType", channelType, "error", err)
		return fmt.Errorf("store channel config. %w", err)
	}
	return nil
}

func (s *ChannelConfigService) Load(ctx context.Context, channelType core.ChannelType, appKey string) (*core.ChannelConfig, error) {
	bucket := s.newBucket(appKey)
	resp := s.rds.HGet(ctx, bucket, string(channelType))
	if err := resp.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		slog.Error("redis load channel config bucket error.", "appKey", appKey, "channelType", channelType, "error", resp.Err())
		return nil, fmt.Errorf("load channel entity %w", resp.Err())
	}
	if resp.Val() != "" {
		var config = &core.ChannelConfig{}
		err := sonic.ConfigFastest.Unmarshal([]byte(resp.Val()), config)
		if err != nil {
			slog.Error("unmarshal channel config error.", "appKey", appKey, "value", resp.Val(), "error", resp.Err())
			return nil, fmt.Errorf("unmarahsl channel config. %w", err)
		}
		return config, nil
	}
	return nil, nil
}

func (s *ChannelConfigService) Exists(ctx context.Context, channelType core.ChannelType, appKey string) bool {
	bucket := s.newBucket(appKey)
	ret := s.rds.HExists(ctx, bucket, string(channelType))
	if ret.Err() != nil {
		slog.Error("redis channel config exists error.", "appKey", appKey, "channelType", channelType, "error", ret.Err())
		return false
	}
	return ret.Val()
}

func (s *ChannelConfigService) newBucket(appKey string) string {
	return fmt.Sprintf("onepush::core::channels::%s", appKey)
}
