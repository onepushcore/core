package core

type Message struct {
	ID      string            `json:"id"`      // 消息ID
	AppKey  string            `json:"app_key"` // 所属应用
	Format  string            `json:"format"`  // 格式化类型
	Title   string            `json:"title"`   // 标题
	Content string            `json:"content"` // 消息内容
	Args    map[string]string `json:"args"`    // 附加参数
}
