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

type UserService struct{}

var userRepo repository.UserRepo

func (u *UserService) Register(ctx context.Context, req *app.RequestContext) {
	username := req.Query("username")
	password := req.Query("password")

	user, err := userRepo.FindUserByUserName(username)
	if err != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	if user != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "用户名已存在",
		})
		return
	}

	userId, err := userRepo.CreateUser(&model.User{
		Name:     username,
		Password: password,
	})
	if err != nil {
		req.JSON(http.StatusOK, dto.BaseResp{
			StatusCode: 0,
			StatusMsg:  "注册失败",
		})
		return
	}
	token, err := tools.GenerateToken(strconv.FormatInt(userId, 10))
	if err != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserId:     userId,
		Token:      token,
	})
}

func (u *UserService) Login(ctx context.Context, req *app.RequestContext) {
	username := req.Query("username")
	password := req.Query("password")

	user, err := userRepo.FindUserByUserName(username)
	if err != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	// 无此用户
	if user == nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}
	// 密码解密
	err = user.DesPassword()
	if err != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	if user.Password != password {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}
	token, err := tools.GenerateToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	req.JSON(http.StatusOK, dto.LoginAndRegisterResp{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     user.ID,
		Token:      token,
	})
}

func (u *UserService) UserInfo(ctx context.Context, req *app.RequestContext) {
	userId := req.Query("user_id")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		req.JSON(http.StatusOK, dto.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	user, err := userRepo.GetUserInfo(id)
	if err != nil {
		req.JSON(http.StatusOK, dto.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		})
		return
	}
	// model -> dto
	userDto := &dto.User{
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		ID:              user.ID,
		IsFollow:        false,
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
	req.JSON(http.StatusOK, dto.UserInfoResp{
		StatusCode: 0,
		StatusMsg:  "拉取用户信息",
		User:       userDto,
	})
}
