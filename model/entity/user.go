package entity

import (
	"tiktok/tools"
)

type User struct {
	Model
	Name            string
	Password        string
	FollowCount     int64
	FollowerCount   int64
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

func (user *User) SetId() {
	user.ID = tools.NewSnowFlake(0).GenSnowID()
}
func (user *User) EncryptPassword() (err error) {
	password, err := tools.EncryptByAes([]byte(user.Password))
	if err != nil {
		return err
	}
	user.Password = password
	return nil
}
func (user *User) DesPassword() (err error) {
	password, err := tools.DecryptByAes(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return nil
}
func (user *User) SetSignature() error {
	signature, err := tools.GeneratePersonalSignature()
	if err != nil {
		return err
	}
	user.Signature = signature
	return nil
}
