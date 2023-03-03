package resp

import (
	"tiktok/model/dto"
)

type BaseResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type LoginAndRegisterResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}
type UserInfoResp struct {
	StatusCode int64     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string    `json:"status_msg"`  // 返回状态描述
	User       *dto.User `json:"user"`        // 用户信息
}
type FeedResp struct {
	NextTime   *int64      `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string     `json:"status_msg"`  // 返回状态描述
	VideoList  []dto.Video `json:"video_list"`  // 视频列表
}
