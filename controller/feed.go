package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
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
	if len(token) == 0 || token == "0" {
		isLogin = false
	}
	// 判断是否首次加载feed流
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
	// 视频库无剩余视频 从头加载
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
	ch1 := make(chan dto.Video, 10)
	ch2 := make(chan dto.Video, 10)
	ch3 := make(chan dto.Video, 10)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for _, video := range videoList {
			item := dto.Video{
				Author:        dto.User{ID: video.Author, IsFollow: false},
				CommentCount:  video.CommentCount,
				PlayURL:       video.PlayUrl,
				CoverURL:      video.CoverUrl,
				FavoriteCount: video.FavoriteCount,
				ID:            video.ID,
				IsFavorite:    false,
				Title:         video.Title,
			}
			ch1 <- item
		}
		close(ch1)
	}()

	go func() {
		for video := range ch1 {
			author, _ := userRepo.FindUserByUserId(video.Author.ID)
			// TODO
			video.Author.Name = author.Name
			ch2 <- video
		}
		close(ch2)
	}()

	go func() {
		for video := range ch2 {
			if isLogin {
				isFollow, _ := feedService.relationRepo.CheckFollow(userId, video.Author.ID)
				video.Author.IsFollow = isFollow
			}
			ch3 <- video
		}
		close(ch3)
	}()

	go func() {
		for video := range ch3 {
			if isLogin {
				isFavorite, _ := feedService.favoriteRepo.CheckFavorite(userId, video.ID)
				video.IsFavorite = isFavorite
			}
			videosList = append(videosList, video)
		}
		wg.Done()
	}()
	wg.Wait()

	if len(videosList) == 0 {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "视频库无视频",
		})
		return
	}

	req.JSON(http.StatusOK, &dto.FeedResp{
		NextTime:   videoList[len(videosList)-1].CreatedAt.Unix(),
		StatusCode: 0,
		StatusMsg:  "拉取视频成功",
		VideoList:  videosList,
	})
}
