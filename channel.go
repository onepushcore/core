package core

type ChannelType string

const (
	ChannelTypeDingTalkApp     ChannelType = "dingtalk-app"
	ChannelTypeDingTalkRobot   ChannelType = "dingtalk-robot"
	ChannelTypeWechatWorkRobot ChannelType = "wechatwork-robot"
	ChannelTypeWechatWorkApp   ChannelType = "wechatwork-app"
	ChannelTypeFeishuApp       ChannelType = "feishu-app"
	ChannelTypeFeishuRobot     ChannelType = "feishu-robot"
)

type ChannelConfig struct {
	Type       ChannelType             `json:"type"`        // 渠道类型
	Name       string                  `json:"name"`        // 渠道名称
	AppKey     string                  `json:"app_key"`     // 所属AppKey
	Enabled    bool                    `json:"enabled"`     // 是否启用渠道
	SendToken  string                  `json:"send_token"`  // 发送消息的访问令牌
	SendSecret string                  `json:"send_secret"` // 发送消息的密钥
	Features   map[FeatureType]Feature `json:"features"`    // 特性列表
}
