package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
	"tiktok/router"
	"tiktok/serverInit"
)

func main() {
	serverInit.InitDatabase()
	hertz := server.Default()
	router.RegisterRouter(hertz)
	logrus.SetReportCaller(true)
	hertz.Spin()
}
