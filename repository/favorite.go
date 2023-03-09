package repository

import (
	"gorm.io/gorm"
	"tiktok/global"
	"tiktok/model"
)

type FavoriteRepo struct{}

func (favoriteRepo *FavoriteRepo) FavoriteAction(favorite *model.Favorite) (err error) {

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var f *model.Favorite
		// 是否有点赞记录
		err = tx.Find(&f, model.Favorite{
			UserId:  favorite.UserId,
			VideoId: favorite.VideoId,
		}).Error
		if err != nil {
			return err
		}
		// 点赞记录更新
		if f.ID == 0 {
			favorite.GenerateID()
			err = tx.Create(favorite).Error
		} else if f.IsFavorite == favorite.IsFavorite {
			return nil
		} else {
			err = tx.Model(favorite).Where(&model.Favorite{UserId: favorite.UserId, VideoId: favorite.VideoId}).Update("is_favorite", favorite.IsFavorite).Error
		}
		if err != nil {
			return err
		}
		videoExpr := "favorite_count - 1"
		userExpr := "favorite_count - 1"
		authorExpr := "total_favorited - 1"
		if favorite.IsFavorite {
			videoExpr = "favorite_count + 1"
			userExpr = "favorite_count + 1"
			authorExpr = "total_favorited + 1"
		}
		// 视频点赞数更新
		var video *model.Video
		err = tx.Model(&video).Where(&model.Model{ID: favorite.VideoId}).Update("favorite_count", gorm.Expr(videoExpr)).Find(&video).Error
		if err != nil {
			return err
		}
		// 用户喜欢数更新
		var user *model.User
		err = tx.Model(user).Where(&model.Model{ID: favorite.UserId}).Update("favorite_count", gorm.Expr(userExpr)).Error
		if err != nil {
			return err
		}
		// 作者获赞数更新
		err = tx.Model(&user).Where(&model.Model{ID: video.Author}).Update("total_favorited", gorm.Expr(authorExpr)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (favoriteRepo *FavoriteRepo) GetUserFavoriteList(userId int64) (videoList []*model.Video, err error) {
	var favorite model.Favorite
	err = global.DB.
		Model(favorite).
		Select("videos.favorite_count", "videos.cover_url").
		Joins("LEFT JOIN videos ON favorites.video_id = videos.id").
		Where(&model.Favorite{UserId: userId}).
		Scan(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (favoriteRepo *FavoriteRepo) CheckFavorite(userId, videoId int64) (bool, error) {
	favorite := model.Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: true,
	}
	var count int64
	err := global.DB.Model(favorite).Where(favorite).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
