package test

import (
	"log"
	"testing"
	"tiktok/model/entity"
	"tiktok/repository"
	"tiktok/serverInit"
	"time"
)

func TestMessageCreate(t *testing.T) {
	serverInit.InitDatabase()
	r := new(repository.MessageRepo)
	milli := time.Now().UnixMilli()
	err := r.SendMessage(&entity.Message{
		ToUserId:   1,
		FromUserId: 2,
		Content:    "123",
		CreateTime: milli,
	})
	if err != nil {
		log.Fatal(err)
	}

}
