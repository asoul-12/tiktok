package repo

import (
	"github.com/sirupsen/logrus"
	"tiktok/model"
	"tiktok/model/dto"
)

type UserRepo struct{}

func (userRepo *UserRepo) CreateUser(user *model.User) (bool, int64) {
	user.SetId()
	user.EncryptPassword()
	user.SetSignature()
	err := baseRepo.Create(user)
	if err != nil {
		logrus.Error(err)
		return false, -1
	}
	return true, user.ID
}
func (userRepo *UserRepo) FindUserByUserName(username string) *model.User {
	var user *model.User
	err := baseRepo.First(&user, model.User{Name: username})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return user
}

//func (userRepo *UserRepo) FindUserByUserId(userId int64) *model.User {
//	var user *model.User
//	err := baseRepo.First(&user, model.User{ID: userId})
//	if err != nil {
//		logger.Error(err)
//		return nil
//	}
//	return user
//}

func (userRepo *UserRepo) GetUserInfo(userId int64) (user *dto.User, err error) {
	err = baseRepo.First(&user, dto.User{ID: userId})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}
