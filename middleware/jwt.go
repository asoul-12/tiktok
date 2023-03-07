package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/tools"
)

func JWT(ctx context.Context, req *app.RequestContext) {
	token := req.Query("token")
	if len(token) == 0 {
		token = req.PostForm("token")
	}
	if len(token) == 0 {
		req.Abort()
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	claims, err := tools.ParseToken(token)
	if err != nil {
		req.Abort()
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	userId := claims.Audience[0]
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		req.Abort()
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	req.Set("userId", id)
	req.Next(ctx)
}
