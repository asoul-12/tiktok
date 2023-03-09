package test

import (
	"fmt"
	"testing"
	"tiktok/model/entity"
	"tiktok/repository"
	"tiktok/serverInit"
)

func TestGetUserPublishVideoList(t *testing.T) {
	serverInit.InitDatabase()
	i := int64(7031070868872626176)
	repo := repository.VideoRepo{}
	list, _ := repo.GetUserPublishList(i)
	fmt.Println(list)
}

func TestFavoriteAction(t *testing.T) {
	serverInit.InitDatabaseTest()
	repo := repository.FavoriteRepo{}
	err := repo.FavoriteAction(&entity.Favorite{
		UserId:     7031070868872626176,
		VideoId:    7031440251318960128,
		IsFavorite: false,
	})
	fmt.Println(err)
}
