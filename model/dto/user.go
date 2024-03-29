package dto

type User struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
}
type FriendUser struct {
	ID      int64  `json:"id"`
	Avatar  string `json:"avatar"`   // 用户头像
	Name    string `json:"name"`     // 用户名称
	Message string `json:"message"`  // 和该好友的最新聊天消息
	MsgType int    `json:"msgType "` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}
