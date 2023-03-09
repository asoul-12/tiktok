package repository

import (
	"errors"
	"gorm.io/gorm"
	"tiktok/model"
)

type UserRepo struct{}

func (userRepo *UserRepo) CreateUser(user *model.User) (int64, error) {
	user.SetId()
	user.EncryptPassword()
	user.SetSignature()
	err := baseRepo.Create(user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (userRepo *UserRepo) FindUserByUserName(username string) (*model.User, error) {
	var user *model.User
	err := baseRepo.Take(&user, &model.User{Name: username})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) FindUserByUserId(userId int64) (*model.User, error) {
	var user *model.User
	err := baseRepo.Take(&user, model.Model{
		ID: userId,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) GetUserInfo(userId int64) (*model.User, error) {
	var user *model.User
	err := baseRepo.Take(&user, model.Model{ID: userId})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
