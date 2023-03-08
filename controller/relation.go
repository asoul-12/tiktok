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

type RelationService struct {
	relationRepo repository.RelationRepo
}

func (relationService *RelationService) FollowList(ctx context.Context, req *app.RequestContext) {
	id := req.Query("user_id")
	tokenId := req.GetInt64("userId")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户id不存在",
		})
		return
	}
	// repo
	list, err := relationService.relationRepo.FollowList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	// model -> dto
	var followList []dto.User
	for _, v := range list {
		flag, _ := relationService.relationRepo.CheckFollow(tokenId, v.ID)
		user := dto.User{Avatar: v.Avatar, Name: v.Name, IsFollow: flag}
		followList = append(followList, user)
	}
	req.JSON(http.StatusOK, dto.FollowListResp{
		StatusCode: 0,
		StatusMsg:  "获取关注列表",
		UserList:   followList,
	})
}

func (relationService *RelationService) FollowAction(ctx context.Context, req *app.RequestContext) {
	toUser := req.Query("to_user_id")
	actionType := req.Query("action_type")
	userId := req.GetInt64("userId")
	var isFollow bool
	// 参数校验
	switch actionType {
	case "1":
		isFollow = true
	case "2":
		isFollow = false
	default:
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "不支持的操作",
		})
		return
	}
	toUserId, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "关注用户不存在",
		})
		return
	}
	// repo
	err = relationService.relationRepo.Follow(&model.Follow{
		UserId:   userId,
		FollowId: toUserId,
		IsFollow: isFollow,
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "关注失败",
		})
		return
	}
	req.JSON(http.StatusOK, dto.BaseResp{
		StatusCode: 0,
		StatusMsg:  "关注成功",
	})
}

func (relationService *RelationService) FollowerList(ctx context.Context, req *app.RequestContext) {
	id := req.Query("user_id")
	tokenId := req.GetInt64("userId")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户id不存在",
		})
		return
	}
	list, err := relationService.relationRepo.FollowerList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "用户id不存在",
		})
		return
	}
	var followerList []dto.User
	for _, v := range list {
		follow, _ := relationService.relationRepo.CheckFollow(tokenId, v.ID)
		user := dto.User{Name: v.Name, Avatar: v.Avatar, IsFollow: follow}
		followerList = append(followerList, user)
	}
	req.JSON(http.StatusOK, dto.FollowListResp{
		StatusCode: 0,
		StatusMsg:  "拉取粉丝列表",
		UserList:   followerList,
	})

}

func (relationService *RelationService) FriendList(ctx context.Context, req *app.RequestContext) {
	userId := req.GetInt64("userId")
	list, err := relationService.relationRepo.FriendList(userId)
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	var friendList []dto.FriendUser
	for _, user := range list {
		item := dto.FriendUser{
			Avatar:  user.Avatar,
			Name:    user.Name,
			Message: "todo",
			MsgType: "0",
		}
		friendList = append(friendList, item)
	}
	req.JSON(http.StatusOK, dto.FriendListResp{
		StatusCode: 0,
		StatusMsg:  "拉取好友列表",
		UserList:   friendList,
	})
}
