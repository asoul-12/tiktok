package repository

import (
	"gorm.io/gorm"
	"tiktok/global"
	"tiktok/model"
)

type RelationRepo struct{}

func (relationRepo *RelationRepo) Follow(follow *model.Follow) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var f *model.Follow
		err := tx.Model(follow).Find(&f, model.Follow{UserId: follow.UserId, FollowId: follow.FollowId}).Error
		if err != nil {
			return err
		}
		// 更新follow记录
		if f.ID == 0 {
			follow.GenerateID()
			err = tx.Model(follow).Create(follow).Error
		} else if f.IsFollow == follow.IsFollow {
			return nil
		} else {
			err = tx.Model(follow).Where(model.Follow{UserId: follow.UserId, FollowId: follow.FollowId}).Update("is_follow", follow.IsFollow).Error
		}
		if err != nil {
			return err
		}
		// 更新粉丝数 关注数
		followExpr := "follow_count - 1"
		followerExpr := "follower_count - 1"
		if follow.IsFollow {
			followExpr = "follow_count + 1"
			followerExpr = "follower_count + 1"
		}
		err = tx.Table("users").Where(model.Model{ID: follow.UserId}).Update("follow_count", gorm.Expr(followExpr)).Error
		if err != nil {
			return err
		}
		err = tx.Table("users").Where(model.Model{ID: follow.FollowId}).Update("follower_count", gorm.Expr(followerExpr)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (relationRepo *RelationRepo) FollowList(userId int64) (userList []*model.User, err error) {
	var follow model.Follow
	err = global.DB.
		Model(follow).
		Select("users.avatar", "users.name").
		Joins("LEFT JOIN users on follows.follow_id = users.id").
		Where(model.Follow{UserId: userId, IsFollow: true}).
		Find(&userList).Error
	if err != nil {
		return nil, err
	}

	return userList, err
}

func (relationRepo *RelationRepo) CheckFollow(userId, targetId int64) (bool, error) {
	var follow *model.Follow
	var count int64
	err := global.DB.Model(follow).Where(&model.Follow{UserId: userId, FollowId: targetId, IsFollow: true}).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}
