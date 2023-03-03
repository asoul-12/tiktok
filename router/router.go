package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"tiktok/handler"
	"time"
)

func RegisterRouter(hertz *server.Hertz) {
	hertz.Use(func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		end := time.Now()
		latency := end.Sub(start).Microseconds
		hlog.CtxTracef(ctx, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s query=%s",
			c.Response.StatusCode(), latency,
			c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host(),
			c.Request.QueryString())
	})
	rootRouter := hertz.Group("/douyin")
	var userService handler.UserService
	var feedService handler.FeedService
	rootRouter.POST("/user/register/", userService.Register)
	rootRouter.POST("/user/login/", userService.Login)
	rootRouter.GET("/user/", userService.UserInfo)
	rootRouter.GET("/feed", feedService.Feed)
	rootRouter.POST("/publish/action", feedService.Publish)
}
