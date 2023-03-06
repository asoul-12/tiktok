package main

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
	"os"
	"tiktok/repository"
	"tiktok/router"
	"tiktok/serverInit"
	"time"
)

func main() {
	serverInit.InitDatabase()
	videoRepo := repository.VideoRepo{}
	videoList := videoRepo.GetFeedList(time.Now().Unix())
	fmt.Println(videoList)
	hertz := server.Default()
	dir, _ := os.Getwd()
	hertz.Static("/assets", dir)
	router.RegisterRouter(hertz)
	logrus.SetReportCaller(true)
	hertz.Spin()
}
