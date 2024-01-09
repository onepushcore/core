package core

type FeatureType string

type FeatureConfig struct {
	AppKey  string         `json:"app_key"`      // 应用标识
	Channel ChannelType    `json:"channel_type"` // 通道类型
	Type    FeatureType    `json:"type"`         // 特性类型
	Enabled bool           `json:"enabled"`      // 是否启用
	Options map[string]any `json:"options"`      // 配置参数
}
