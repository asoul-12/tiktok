package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/model"
	"tiktok/model/dto/resp"
	"tiktok/repo"
	"tiktok/tools"
)

type UserService struct{}

var userRepo repo.UserRepo

func (u *UserService) Register(ctx context.Context, req *app.RequestContext) {
	username := req.Query("username")
	password := req.Query("password")

	user := userRepo.FindUserByUserName(username)
	if user != nil {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "用户名已存在",
		})
		return
	}

	isSuccess, id := userRepo.CreateUser(&model.User{
		Name:     username,
		Password: password,
	})
	if isSuccess {
		token, err := tools.GenerateToken(strconv.FormatInt(id, 10))
		if err == nil {
			req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
				StatusCode: 0,
				StatusMsg:  "注册成功",
				UserId:     id,
				Token:      token,
			})
			return
		}
	}

	req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
		StatusCode: 1,
		StatusMsg:  "网络出错了",
	})
}

func (u *UserService) Login(ctx context.Context, req *app.RequestContext) {
	username := req.Query("username")
	password := req.Query("password")

	user := userRepo.FindUserByUserName(username)
	// 无此用户
	if user == nil {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}
	// 密码解密
	err := user.DesPassword()
	if err != nil {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	if user.Password != password {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}
	token, err := tools.GenerateToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     user.ID,
		Token:      token,
	})
}

func (u *UserService) UserInfo(ctx context.Context, req *app.RequestContext) {
	userId, err := strconv.ParseInt(req.Query("user_id"), 10, 64)
	if err != nil {
		logrus.Error(err)
		req.JSON(http.StatusOK, resp.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	user, err := userRepo.GetUserInfo(userId)
	if err != nil {
		req.JSON(http.StatusOK, resp.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	req.JSON(http.StatusOK, resp.UserInfoResp{
		StatusCode: 0,
		StatusMsg:  "拉取用户信息",
		User:       user,
	})
}
