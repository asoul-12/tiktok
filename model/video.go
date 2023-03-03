package model

type Video struct {
	Model
	Author        int64
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int
	CommentCount  int
	IsFavorite    bool
	Title         string
}
