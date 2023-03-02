package handler

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
	"strconv"
	"tiktok/api"
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

	id := tools.NewSnowFlake(0).GenSnowID()
	password, err := tools.EncryptByAes([]byte(password))
	if err != nil {
		logger.Error(err)
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络出错了",
		})
		return
	}
	isSuccess := userRepo.CreateUser(&model.User{
		ID:       id,
		Name:     username,
		Password: password,
		//Avatar:          "",
		//BackgroundImage: "",
		Signature: api.GeneratePersonalSignature(),
	})
	if isSuccess {
		token, err := tools.GenerateToken(username)
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
	logger.Error(err)
	req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
		StatusCode: 1,
		StatusMsg:  "网络出错了",
	})
}

func (u *UserService) Login(ctx context.Context, req *app.RequestContext) {
	username := req.Query("username")
	password := req.Query("password")

	user := userRepo.FindUserByUserName(username)
	dbPwd, err := tools.DecryptByAes(user.Password)
	pwd := string(dbPwd)
	if err != nil {
		log.Print(err)
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	if pwd != password {
		req.JSON(http.StatusOK, resp.LoginAndRegisterResp{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}

	token, err := tools.GenerateToken(username)
	if err != nil {
		log.Print(err)
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
	println(req.Query("user_id"))
	userId, err := strconv.ParseInt(req.Query("user_id"), 10, 64)
	token := req.Query("token")
	if err != nil {
		log.Print(err)
		req.JSON(http.StatusOK, resp.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "网络错误",
		})
		return
	}
	_, err = tools.ParseToken(token)
	if err != nil {
		log.Print(err)
		req.JSON(http.StatusOK, resp.UserInfoResp{
			StatusCode: 1,
			StatusMsg:  "token过期",
		})
		return
	}

	user := userRepo.FindUserByUserId(userId)
	log.Print(user)
	req.JSON(http.StatusOK, resp.UserInfoResp{
		StatusCode: 0,
		StatusMsg:  "拉取用户信息",
		User:       user,
	})

}
