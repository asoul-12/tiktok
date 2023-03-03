package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/model/dto/resp"
	"tiktok/tools"
)

type FeedService struct{}

func (feedService *FeedService) Feed(ctx context.Context, req *app.RequestContext) {

}

func (feedService *FeedService) Publish(ctx context.Context, req *app.RequestContext) {
	file, err := req.FormFile("data")
	file1 := file
	title := req.PostForm("title")
	token := req.PostForm("token")
	if err != nil {
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	claims, err := tools.ParseToken(token)
	if err != nil {
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	userId, err := strconv.ParseInt(claims.Audience[0], 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	byteData := make([]byte, file.Size)
	open, err := file.Open()
	if err != nil {
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	open.Read(byteData)
	contentType := http.DetectContentType(byteData)
	if contentType != "video/mp4" {
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "文件格式有误",
		})
		return
	}
	println(userId, title, contentType)
	err = req.SaveUploadedFile(file1, "D:\\code-02\\Go_workspace\\src\\tiktok\\assets\\"+title+file1.Filename)
	if err != nil {
		logrus.Error(err)
		req.JSON(http.StatusOK, resp.BaseResp{
			StatusCode: 1,
			StatusMsg:  "上传失败",
		})
		return
	}
	req.JSON(http.StatusOK, resp.BaseResp{
		StatusCode: 0,
		StatusMsg:  "上传成功",
	})
}
