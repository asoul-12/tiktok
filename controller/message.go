package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
	"tiktok/model/dto"
	"tiktok/model/entity"
	"tiktok/model/query"
	"tiktok/repository"
	"time"
)

type MessagesService struct {
	messageRepo  repository.MessageRepo
	relationRepo repository.RelationRepo
}

func (messageService *MessagesService) SendMessage(ctx context.Context, req *app.RequestContext) {
	fromUserId := req.GetInt64("userId")
	toUser := req.Query("to_user_id")
	actionType := req.Query("action_type")
	content := req.Query("content")
	// 参数校验
	toUserId, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "对方id不存在",
		})
		return
	}
	// 是否互关
	isFriend, err := messageService.relationRepo.CheckFriend(fromUserId, toUserId)
	if actionType != "1" || !isFriend || err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "不支持的操作",
		})
		return
	}
	// repo
	err = messageService.messageRepo.SendMessage(&entity.Message{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
		CreateTime: time.Now().UnixMilli(),
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "发送成功",
	})
}

func (messageService *MessagesService) MessageList(ctx context.Context, req *app.RequestContext) {
	preMsgTime := req.Query("pre_msg_time")
	toUser := req.Query("to_user_id")
	toUserId, err := strconv.ParseInt(toUser, 10, 64)
	userId := req.GetInt64("userId")
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "对方id不存在",
		})
		return
	}
	// 判断是否第一次拉取信息
	reqTime, err := strconv.ParseInt(preMsgTime, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "请求时间有误",
		})
		return
	}
	// repo
	list, err := messageService.messageRepo.MessageList(query.MessageQuery{
		FromUserId: userId,
		ToUserId:   toUserId,
		CreateTime: reqTime,
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错",
		})
		return
	}
	// model -> dto
	var messageList []dto.Message
	for _, message := range list {
		item := dto.Message{
			Content:    message.Content,
			CreateTime: message.CreateTime,
			FromUserID: message.FromUserId,
			ID:         message.ID,
			ToUserID:   message.ToUserId,
		}
		messageList = append(messageList, item)
	}
	req.JSON(http.StatusOK, dto.MessageResp{
		MessageList: messageList,
		StatusCode:  0,
		StatusMsg:   "拉取信息记录",
	})
}
