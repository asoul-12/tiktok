package controller

import (
	"context"
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

type PublishService struct {
	videoRepo repository.VideoRepo
}

var ffmpeg = tools.Bind{
	FFMpeg:         "D:\\tools\\ffmpeg\\bin\\ffmpeg.exe",
	FFProbe:        "D:\\tools\\ffmpeg\\bin\\ffprobe.exe",
	CommandTimeout: 5000,
}

func (publishService *PublishService) Publish(ctx context.Context, req *app.RequestContext) {
	file, err := req.FormFile("data")
	title := req.PostForm("title")
	userId := req.GetInt64("userId")
	// 视频格式校验
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}

	videoExt := path.Ext(file.Filename)
	jpgExt := ".jpg"
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
	playUrl := "D:\\code-02\\Go_workspace\\src\\tiktok\\assets\\play\\" + fileName + videoExt
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
	coverUrl := "D:\\code-02\\Go_workspace\\src\\tiktok\\assets\\cover\\" + fileName + jpgExt
	err = ffmpeg.Thumbnail(playUrl, coverUrl, 0, true)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "生成缩略图失败",
		})
		return
	}
	// repository
	playUrl = "http://192.168.31.187:8888/assets/play/" + fileName + videoExt
	coverUrl = "http://192.168.31.187:8888/assets/cover/" + fileName + jpgExt
	err = publishService.videoRepo.CreateVideo(&model.Video{
		Author:        userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "上传失败",
		})
		return
	}
	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "上传成功",
	})
}

func (publishService *PublishService) PublishList(ctx context.Context, req *app.RequestContext) {
	userId := req.Query("user_id")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	// repo
	videoList, err := publishService.videoRepo.GetUserPublishList(id)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "加载作品列表失败",
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
		StatusMsg:  "拉取用户作品列表",
		VideoList:  videosList,
	})
}
