package test

import (
	"fmt"
	"log"
	"testing"
	"tiktok/repository"
	"tiktok/serverInit"
	"time"
)

func TestUser(t *testing.T) {
	serverInit.InitDatabase()
	userRepo := repository.UserRepo{}
	name, err := userRepo.FindUserByUserName("asoul")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}

func TestGetVideoList(t *testing.T) {
	serverInit.InitDatabase()
	videoRepo := repository.VideoRepo{}
	list, _ := videoRepo.GetFeedList(time.Now().Unix())
	videoList := list
	fmt.Println(videoList)
}
