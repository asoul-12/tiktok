package entity

type Message struct {
	Model
	ToUserId   int64
	FromUserId int64
	Content    string
	CreateTime int64
}
