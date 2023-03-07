package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
	"os"
	"tiktok/router"
	"tiktok/serverInit"
)

func main() {
	serverInit.InitDatabase()
	hertz := server.Default()
	dir, _ := os.Getwd()
	hertz.Static("/assets", dir)
	router.RegisterRouter(hertz)
	logrus.SetReportCaller(true)
	hertz.Spin()
}
