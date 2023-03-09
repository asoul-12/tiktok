package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/model/entity"
	"tiktok/repository"
	"time"
)

type CommentService struct {
	commentRepo repository.CommentRepo
	userRepo    repository.UserRepo
}

func (commentService *CommentService) CommentAction(ctx context.Context, req *app.RequestContext) {
	userId := req.GetInt64("userId")
	videoID := req.Query("video_id")
	actionType := req.Query("action_type")
	var commentText string
	var commentId int64
	var err error
	var comment entity.Comment
	// 参数校验
	videoId, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "未知视频id",
		})
		return
	}
	switch actionType {
	case "1":
		commentText = req.Query("comment_text")
		comment = entity.Comment{
			UserId:     userId,
			VideoId:    videoId,
			Content:    commentText,
			CreateDate: time.Now().Unix(),
		}
		err = commentService.commentRepo.AddComment(&comment)
	case "2":
		commentId, err = strconv.ParseInt(req.Query("comment_id"), 10, 64)
		if err != nil {
			req.JSON(http.StatusOK, dto.BaseResp{
				StatusCode: 1,
				StatusMsg:  "指定评论不存在",
			})
			return
		}
		err = commentService.commentRepo.DelComment(commentId)
	default:
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "不支持的操作",
		})
		return
	}
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	var commentDto dto.Comment
	if actionType == "1" {
		commentDto.ID = comment.ID
		commentDto.Content = comment.Content
		commentDto.CreateDate = time.Unix(comment.CreateDate, 0).Format("2006-01-02 15:04:05")
		userInfo, _ := userRepo.FindUserByUserId(userId)
		commentDto.User = dto.User{
			Avatar: userInfo.Avatar,
			ID:     userInfo.ID,
			Name:   userInfo.Name,
		}
	}
	req.JSON(http.StatusOK, dto.CommentActionResp{
		StatusCode: 0,
		StatusMsg:  "操作成功",
		Comment:    commentDto,
	})
}

func (commentService *CommentService) CommentList(ctx context.Context, req *app.RequestContext) {
	videoID := req.Query("video_id")
	// 参数校验
	videoId, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "视频不存在",
		})
		return
	}
	// repo
	list, err := commentService.commentRepo.CommentList(videoId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	// model -> dto
	var commentList []dto.Comment
	for _, comment := range list {
		userInfo, _ := userRepo.GetUserInfo(comment.UserId)
		item := dto.Comment{
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CreateDate, 0).Format("2006-01-02 15:04:05"),
			ID:         comment.ID,
			User: dto.User{
				ID:     userInfo.ID,
				Avatar: userInfo.Avatar,
				Name:   userInfo.Name,
			},
		}
		commentList = append(commentList, item)
	}
	req.JSON(http.StatusOK, dto.CommentListResp{
		StatusCode:  0,
		StatusMsg:   "拉取视频评论",
		CommentList: commentList,
	})
}
