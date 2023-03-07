package model

type Follow struct {
	Model
	UserId   int64
	FollowId int64
	IsFollow bool
}
