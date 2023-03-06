package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"tiktok/model"
)

type VideoRepo struct{}

func (videoRepo *VideoRepo) CreateVideo(video *model.Video) error {
	video.GenerateID()
	err := baseRepo.Create(video)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
func (videoRepo *VideoRepo) GetFeedList(timestamp int64) []*model.Video {
	var videoList []*model.Video
	where := fmt.Sprintf("UNIX_TIMESTAMP(created_at) < %d", timestamp)
	err := baseRepo.FindWhereOrderLimit(&videoList, where, "created_at DESC", 5)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return videoList
}
func (videoRepo *VideoRepo) GetUserFavoriteList(userId int64) (videoList []*model.Video, err error) {
	err = baseRepo.Find(&videoList, &model.Video{Author: userId})
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
