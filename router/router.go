package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/controller"
	"tiktok/middleware"
)

var (
	userService     controller.UserService
	feedService     controller.FeedService
	publishService  controller.PublishService
	favoriteService controller.FavoriteService
	relationService controller.RelationService
)

func RegisterRouter(hertz *server.Hertz) {
	hertz.Use(middleware.Log)
	rootRouter := hertz.Group("/douyin")
	authRouter := rootRouter.Group("/", middleware.JWT)
	// user
	rootRouter.POST("/user/register/", userService.Register)
	rootRouter.POST("/user/login/", userService.Login)
	authRouter.GET("/user/", userService.UserInfo)
	// feed
	rootRouter.GET("/feed/", feedService.Feed)
	// publish
	authRouter.GET("/publish/list/", publishService.PublishList)
	authRouter.POST("/publish/action/", publishService.Publish)
	// favorite
	authRouter.POST("/favorite/action/", favoriteService.Action)
	authRouter.GET("/favorite/list/", favoriteService.List)
	// relation
	authRouter.GET("/relation/follow/list/", relationService.FollowList)
	authRouter.POST("/relation/action/", relationService.FollowAction)
	authRouter.GET("/relation/follower/list/", relationService.FollowerList)
	authRouter.GET("/relation/friend/list/", relationService.FriendList)
	// message

}
