package dto

type Message struct {
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  `json:"from_user_id"` // 消息发送者id
	ID         int64  `json:"id"`           // 消息id
	ToUserID   int64  `json:"to_user_id"`   // 消息接收者id
}
