package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/handler"
)

func RegisterRouter(hertz *server.Hertz) {
	rootRouter := hertz.Group("/douyin")
	var userService handler.UserService
	rootRouter.POST("/user/register/", userService.Register)
	rootRouter.POST("/user/login/", userService.Login)
	rootRouter.GET("/user/", userService.UserInfo)
}
