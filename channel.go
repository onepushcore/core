package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

const (
	ChannelTypeDingTalkApp     = "dingtalk-app"
	ChannelTypeDingTalkRobot   = "dingtalk-robot"
	ChannelTypeWechatWorkRobot = "wechatwork-robot"
	ChannelTypeWechatWorkApp   = "wechatwork-app"
	ChannelTypeFeishuApp       = "feishu-app"
	ChannelTypeFeishuRobot     = "feishu-robot"
)

type ChannelEntity struct {
	Type       string         `json:"type"`        // 渠道类型
	AppKey     string         `json:"app_key"`     // 所属AppKey
	Enabled    bool           `json:"enabled"`     // 是否启用渠道
	SendToken  string         `json:"send_token"`  // 发送消息的访问令牌
	SendSecret string         `json:"send_secret"` // 发送消息的密钥
	Attrs      map[string]any `json:"attrs"`       // 属性列表
}

////

type ChannelService struct {
	bucket *BucketService[ChannelEntity]
}

func NewChannelService(rds redis.Cmdable) *ChannelService {
	return &ChannelService{
		bucket: &BucketService[ChannelEntity]{
			Rds: rds,
			NewBucket: func(owner string) string {
				return "onepush::core::channels::" + owner
			},
			NewValue: func() ChannelEntity {
				return ChannelEntity{}
			},
		},
	}
}

func (s *ChannelService) List(ctx context.Context, appKey string) (map[string]ChannelEntity, error) {
	dmap, err := s.bucket.List(ctx, appKey)
	if err != nil {
		slog.Error("redis list channel entity error.", "appKey", appKey, "error", err)
		return nil, err
	}
	return dmap, nil
}

func (s *ChannelService) Store(ctx context.Context, channel ChannelEntity) error {
	err := s.bucket.Store(ctx, channel.AppKey, channel.Type, channel)
	if err != nil {
		slog.Error("redis store channel entity error.", "appKey", channel.AppKey, "channelType", channel.Type, "error", err)
		return err
	}
	return nil
}

func (s *ChannelService) Load(ctx context.Context, appKey string, channelType string) (*ChannelEntity, error) {
	value, err := s.bucket.Load(ctx, appKey, channelType)
	if err != nil {
		slog.Error("redis load channel entity error.", "appKey", appKey, "channelType", channelType, "error", err)
		return nil, err
	}
	return value, nil
}

func (s *ChannelService) Exists(ctx context.Context, appKey string, channelType string) bool {
	exists, err := s.bucket.Exists(ctx, appKey, channelType)
	if err != nil {
		slog.Error("redis exists channel entity error.", "appKey", appKey, "channelType", channelType, "error", err)
		return false
	}
	return exists
}

func (s *ChannelService) Remove(ctx context.Context, appKey string, channelType string) error {
	err := s.bucket.Remove(ctx, appKey, channelType)
	if err != nil {
		slog.Error("redis remove channel entity error.", "appKey", appKey, "channelType", channelType, "error", err)
		return err
	}
	return nil
}
