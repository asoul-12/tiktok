package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strconv"
	"tiktok/model"
	"tiktok/model/dto"
	"tiktok/repository"
	"tiktok/tools"
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
	videoList := feedService.videoRepo.GetFeedList(reqTime)
	if len(videoList) == 0 {
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

func (feedService *FeedService) Publish(ctx context.Context, req *app.RequestContext) {
	file, err := req.FormFile("data")
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	title := req.PostForm("title")
	token := req.PostForm("token")
	// 解析 token
	claims, err := tools.ParseToken(token)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	userId, err := strconv.ParseInt(claims.Audience[0], 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	// 视频校验
	ext := path.Ext(file.Filename)
	byteData := make([]byte, file.Size)
	open, err := file.Open()
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	open.Read(byteData)
	contentType := http.DetectContentType(byteData)
	if contentType != "video/mp4" {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	// 本地存储
	fileName := title + strconv.FormatInt(time.Now().UnixMilli(), 10)
	playUrl := "D:\\code-02\\Go_workspace\\src\\tiktok\\assets\\play\\" + fileName + ext
	err = req.SaveUploadedFile(file, playUrl)
	if err != nil {
		logrus.Error(err)
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "上传失败",
		})
		return
	}
	// 缩略图
	ffmpeg := tools.Bind{
		FFMpeg:         "D:\\tools\\ffmpeg\\bin\\ffmpeg.exe",
		FFProbe:        "D:\\tools\\ffmpeg\\bin\\ffprobe.exe",
		CommandTimeout: 5000,
	}
	coverUrl := "D:\\code-02\\Go_workspace\\src\\tiktok\\assets\\cover\\" + fileName + ".jpg"
	err = ffmpeg.Thumbnail(playUrl, coverUrl, 0, true)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "上传失败",
		})
		return
	}
	// repository
	playUrl = "http://192.168.31.187:8888/assets/play/" + fileName + ext
	coverUrl = "http://192.168.31.187:8888/assets/cover/" + fileName + ".jpg"
	feedService.videoRepo.CreateVideo(&model.Video{
		Author:        userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})

	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "上传成功",
	})
}

func (feedService *FeedService) PublishList(ctx context.Context, req *app.RequestContext) {
	userId, err := strconv.ParseInt(req.Query("user_id"), 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户id无法识别",
		})
		return
	}
	videoList, err := feedService.videoRepo.GetUserPublishList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载作品列表失败",
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
		StatusMsg:  "拉取用户作品列表",
		VideoList:  videosList,
	})
}
