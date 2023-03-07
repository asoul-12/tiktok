package repository

import (
	"fmt"
	"gorm.io/gorm"
	"tiktok/global"
	"tiktok/model"
)

type VideoRepo struct{}

func (videoRepo *VideoRepo) CreateVideo(video *model.Video) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 插入视频记录
		video.GenerateID()
		err := tx.Create(video).Error
		if err != nil {
			return err
		}
		// 用户视频数增加
		var user *model.User
		err = tx.Model(user).Where(&model.Model{ID: video.Author}).Update("work_count", gorm.Expr("work_count + 1")).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
func (videoRepo *VideoRepo) GetFeedList(timestamp int64) ([]*model.Video, error) {
	var videoList []*model.Video
	where := fmt.Sprintf("UNIX_TIMESTAMP(created_at) < %d", timestamp)
	err := baseRepo.FindWhereOrderLimit(&videoList, where, "created_at DESC", 5)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (videoRepo *VideoRepo) GetUserPublishList(userId int64) (videoList []*model.Video, err error) {
	err = baseRepo.Find(&videoList, &model.Video{Author: userId})
	if err != nil {
		return nil, err
	}
	return videoList, nil
}
