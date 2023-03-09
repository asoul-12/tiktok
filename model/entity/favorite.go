package entity

type Favorite struct {
	Model
	UserId     int64
	VideoId    int64
	IsFavorite bool
}
