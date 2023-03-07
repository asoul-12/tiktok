package repository

import "tiktok/model"

type FavoriteRepo struct{}

func (favoriteRepo *FavoriteRepo) FavoriteAction(favorite *model.Favorite) (err error) {
	var f *model.Favorite
	err = baseRepo.Find(&f, model.Favorite{
		UserId:  favorite.UserId,
		VideoId: favorite.VideoId,
	})
	if err != nil {
		return err
	}
	if f.ID == 0 {
		favorite.GenerateID()
		err = baseRepo.Create(favorite)
	} else {
		err = baseRepo.update(favorite, &model.Favorite{UserId: favorite.UserId, VideoId: favorite.VideoId}, "is_favorite", favorite.IsFavorite)
	}
	if err != nil {
		return err
	}

	return nil
}
