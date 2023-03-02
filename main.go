package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/handler"
	"tiktok/router"
	"tiktok/serverInit"
)

func main() {
	serverInit.InitDatabase()
	hertz := server.Default()
	router.RegisterRouter(hertz)
	var userService handler.UserService
	hertz.GET("/ping", userService.Register)

	hertz.Spin()
}
