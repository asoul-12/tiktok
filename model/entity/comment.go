package entity

type Comment struct {
	Model
	UserId     int64
	VideoId    int64
	Content    string
	CreateDate int64
}
