package test

import (
	"fmt"
	"testing"
	"tiktok/repository"
	"tiktok/serverInit"
	"time"
)

func TestUser(t *testing.T) {
	serverInit.InitDatabase()
	userRepo := repository.UserRepo{}
	name := userRepo.FindUserByUserName("asoul")
	fmt.Println(name)
}

func TestGetVideoList(t *testing.T) {
	serverInit.InitDatabaseTest()
	videoRepo := repository.VideoRepo{}
	list, _ := videoRepo.GetFeedList(time.Now().Unix())
	videoList := list
	fmt.Println(videoList)
}
