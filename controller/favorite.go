package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/repository"
)

type FavoriteService struct {
	videoRepo repository.VideoRepo
}

func (favoriteService *FavoriteService) Action(ctx context.Context, req *app.RequestContext) {

}

func (favoriteService *FavoriteService) List(ctx context.Context, req *app.RequestContext) {
	//token := req.Query("token")
	userId, err := strconv.ParseInt(req.Query("user_id"), 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户id无法识别",
		})
		return
	}
	videoList, err := favoriteService.videoRepo.GetUserPublishList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载喜欢列表失败",
		})
		return
	}
	var videosList []dto.Video
	for _, v := range videoList {
		item := dto.Video{
			FavoriteCount: v.FavoriteCount,
			CoverURL:      v.CoverUrl,
		}
		videosList = append(videosList, item)
	}
	req.JSON(http.StatusOK, dto.VideoListResp{
		StatusCode: "0",
		StatusMsg:  "拉取用户喜欢列表",
		VideoList:  videosList,
	})
}
