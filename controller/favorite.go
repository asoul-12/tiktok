package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model"
	"tiktok/model/dto"
	"tiktok/repository"
	"tiktok/tools"
)

type FavoriteService struct {
	videoRepo    repository.VideoRepo
	favoriteRepo repository.FavoriteRepo
}

func (favoriteService *FavoriteService) Action(ctx context.Context, req *app.RequestContext) {
	video := req.Query("video_id")
	actionType := req.Query("action_type")
	token := req.Query("token")
	isFavorite := false
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
	// token 解析
	claims, err := tools.ParseToken(token)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token过期",
		})
		return
	}
	user := claims.Audience[0]
	userId, userErr := strconv.ParseInt(user, 10, 64)
	videoId, VideoErr := strconv.ParseInt(video, 10, 64)
	if userErr != nil || VideoErr != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token过期",
		})
		return
	}
	// repo
	favoriteService.favoriteRepo.FavoriteAction(&model.Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: isFavorite,
	})
	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "点赞",
	})
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
