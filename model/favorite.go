package model

type Favorite struct {
	Model
	UserId     int64
	VideoId    int64
	isFavorite bool
}
