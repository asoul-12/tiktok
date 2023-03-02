package repo

import (
	"github.com/bytedance/gopkg/util/logger"
	"tiktok/model"
)

type UserRepo struct{}

func (userRepo *UserRepo) CreateUser(user *model.User) bool {
	err := baseRepo.Create(user)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}
func (userRepo *UserRepo) FindUserByUserName(username string) *model.User {
	var user *model.User
	err := baseRepo.First(&user, model.User{Name: username})
	if err != nil {
		logger.Error(err)
		return nil
	}
	return user
}

func (userRepo *UserRepo) FindUserByUserId(userId int64) *model.User {
	var user *model.User
	err := baseRepo.First(&user, model.User{ID: userId})
	if err != nil {
		logger.Error(err)
		return nil
	}
	return user
}

func (userRepo *UserRepo) CheckUser(username, password string) (*model.User, error) {
	var user *model.User
	err := baseRepo.First(&user, model.User{Name: username, Password: password})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return user, nil
}
