package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model"
	"tiktok/model/dto"
	"tiktok/repository"
)

type FavoriteService struct {
	videoRepo    repository.VideoRepo
	favoriteRepo repository.FavoriteRepo
}

func (favoriteService *FavoriteService) Action(ctx context.Context, req *app.RequestContext) {
	video := req.Query("video_id")
	actionType := req.Query("action_type")
	userId := req.GetInt64("userId")
	var isFavorite bool
	// 参数校验
	switch actionType {
	case "1":
		isFavorite = true
	case "2":
		isFavorite = false
	default:
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "不支持的操作",
		})
		return
	}
	videoId, VideoErr := strconv.ParseInt(video, 10, 64)
	if VideoErr != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "视频不存在",
		})
		return
	}
	// repo
	err := favoriteService.favoriteRepo.FavoriteAction(&model.Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: isFavorite,
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "点赞失败",
		})
		return
	}
	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "点赞",
	})
}

func (favoriteService *FavoriteService) List(ctx context.Context, req *app.RequestContext) {
	id := req.Query("user_id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		if err != nil {
			req.JSON(http.StatusOK, dto.BaseResp{
				StatusCode: 1,
				StatusMsg:  "用户id不存在",
			})
			return
		}
	}
	// repo
	videoList, err := favoriteService.favoriteRepo.GetUserFavoriteList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载喜欢列表失败",
		})
		return
	}
	// model -> dto
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
