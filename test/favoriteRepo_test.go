package test

import (
	"fmt"
	"testing"
	"tiktok/repository"
	"tiktok/serverInit"
)

func TestGetUserVideoList(t *testing.T) {
	serverInit.InitDatabase()
	i := int64(7031070868872626176)
	repo := repository.VideoRepo{}
	list, _ := repo.GetUserPublishList(i)
	fmt.Println(list)
}
