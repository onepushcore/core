package core

type FeatureType string

const (
	ChannelFeatureTypeTrimSpace FeatureType = "feat:trim-space"
)

type Feature struct {
	Type    FeatureType    `json:"type"`
	Name    string         `json:"name"`    // 特性名称
	Enabled bool           `json:"enabled"` // 是否启用特性
	Options map[string]any `json:"options"` // 配置参数
}
