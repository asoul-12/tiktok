package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/repository"
	"time"
)

type FeedService struct {
	videoRepo repository.VideoRepo
	userRepo  repository.UserRepo
}

func (feedService *FeedService) Feed(ctx context.Context, req *app.RequestContext) {
	token := req.Query("token")
	latestTime := req.Query("latest_time")
	isLogin := true
	reqTime := time.Now().UnixMilli()
	var err error
	// 是否登录以加载点赞状态
	if len(token) == 0 {
		isLogin = false
	}
	fmt.Println(isLogin)
	// 首次加载feed流
	if len(latestTime) != 0 {
		reqTime, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			logrus.Error(err)
			req.JSON(http.StatusOK, dto.BaseResp{
				StatusCode: 1,
				StatusMsg:  "系统时间有误",
			})
		}
		return
	}
	// repo 根据时间戳加载视频
	videoList, err := feedService.videoRepo.GetFeedList(reqTime)
	if len(videoList) == 0 || err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载视频失败，服务器无视频",
		})
		return
	}
	// 拼装dto对象 1.加载该视频的作者信息 （用户是否关注该作者） 2.判断该视频是否被点赞
	var videosList []dto.Video
	for _, v := range videoList {
		author, err := userRepo.FindUserByUserId(v.Author)
		if err != nil {
			req.JSON(http.StatusOK, dto.BaseResp{
				StatusCode: 1,
				StatusMsg:  "加载视频作者信息失败",
			})
			return
		}
		video := dto.Video{
			Author: dto.User{
				Avatar:          author.Avatar,
				BackgroundImage: author.BackgroundImage,
				FavoriteCount:   author.FavoriteCount,
				FollowCount:     author.FollowCount,
				FollowerCount:   author.FollowerCount,
				ID:              author.ID,
				IsFollow:        false,
				Name:            author.Name,
				Signature:       author.Signature,
				TotalFavorited:  author.TotalFavorited,
				WorkCount:       author.WorkCount,
			},
			CommentCount:  0,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			ID:            v.ID,
			IsFavorite:    false,
			PlayURL:       v.PlayUrl,
			Title:         v.Title,
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
