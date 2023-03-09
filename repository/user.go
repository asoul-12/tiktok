package repository

import (
	"errors"
	"gorm.io/gorm"
	"tiktok/model/entity"
)

type UserRepo struct{}

func (userRepo *UserRepo) CreateUser(user *entity.User) (int64, error) {
	user.SetId()
	user.EncryptPassword()
	user.SetSignature()
	err := baseRepo.Create(user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (userRepo *UserRepo) FindUserByUserName(username string) (*entity.User, error) {
	var user *entity.User
	err := baseRepo.Take(&user, &entity.User{Name: username})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) FindUserByUserId(userId int64) (*entity.User, error) {
	var user *entity.User
	err := baseRepo.Take(&user, entity.Model{
		ID: userId,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) GetUserInfo(userId int64) (*entity.User, error) {
	var user *entity.User
	err := baseRepo.Take(&user, entity.Model{ID: userId})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
