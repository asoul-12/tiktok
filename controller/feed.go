package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/repository"
	"time"
)

type FeedService struct {
	videoRepo    repository.VideoRepo
	userRepo     repository.UserRepo
	relationRepo repository.RelationRepo
	favoriteRepo repository.FavoriteRepo
}

func (feedService *FeedService) Feed(ctx context.Context, req *app.RequestContext) {
	token := req.Query("token")
	userId := req.GetInt64("userId")
	latestTime := req.Query("latest_time")
	isLogin := true
	reqTime := time.Now().UnixMilli()
	var err error
	// 是否登录以加载点赞状态
	if len(token) == 0 {
		isLogin = false
	}
	// 首次加载feed流
	if len(latestTime) != 0 && latestTime != "0" {
		reqTime, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			logrus.Error(err)
			req.JSON(http.StatusOK, dto.BaseResp{
				StatusCode: 1,
				StatusMsg:  "系统时间有误",
			})
			return
		}
	}
	// repo 根据时间戳加载视频
	videoList, err := feedService.videoRepo.GetFeedList(reqTime)
	// 视频库无视频 重新加载
	if err == nil && len(videoList) == 0 {
		reqTime = time.Now().UnixMilli()
		videoList, err = feedService.videoRepo.GetFeedList(reqTime)
	}
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载视频失败",
		})
		return
	}

	// 拼装dto对象 1.加载该视频的作者信息 （用户是否关注该作者） 2.判断该视频是否被点赞
	var videosList []dto.Video
	for _, video := range videoList {
		author, _ := userRepo.FindUserByUserId(video.Author)
		// 判断使用者与视频作者关系
		// 判断使用者是否点赞视频
		isFollow := false
		isFavorite := false
		if isLogin {
			isFollow, _ = feedService.relationRepo.CheckFollow(userId, video.Author)
			isFavorite, _ = feedService.favoriteRepo.CheckFavorite(userId, video.ID)
		}
		video := dto.Video{
			Author: dto.User{
				Avatar:          author.Avatar,
				BackgroundImage: author.BackgroundImage,
				FavoriteCount:   author.FavoriteCount,
				FollowCount:     author.FollowCount,
				FollowerCount:   author.FollowerCount,
				ID:              author.ID,
				IsFollow:        isFollow,
				Name:            author.Name,
				Signature:       author.Signature,
				TotalFavorited:  author.TotalFavorited,
				WorkCount:       author.WorkCount,
			},
			CommentCount:  video.CommentCount,
			CoverURL:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			ID:            video.ID,
			IsFavorite:    isFavorite,
			PlayURL:       video.PlayUrl,
			Title:         video.Title,
		}
		videosList = append(videosList, video)
	}
	req.JSON(http.StatusOK, &dto.FeedResp{
		NextTime:   videoList[len(videosList)-1].CreatedAt.Unix(),
		StatusCode: 0,
		StatusMsg:  "拉取视频成功",
		VideoList:  videosList,
	})
}
