package repository

import (
	"gorm.io/gorm"
	"tiktok/global"
	"tiktok/model/entity"
)

type RelationRepo struct{}

func (relationRepo *RelationRepo) Follow(follow *entity.Follow) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var f *entity.Follow
		err := tx.Model(follow).Find(&f, entity.Follow{UserId: follow.UserId, FollowId: follow.FollowId}).Error
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
			err = tx.Model(follow).Where(entity.Follow{UserId: follow.UserId, FollowId: follow.FollowId}).Update("is_follow", follow.IsFollow).Error
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
		err = tx.Table("users").Where(entity.Model{ID: follow.UserId}).Update("follow_count", gorm.Expr(followExpr)).Error
		if err != nil {
			return err
		}
		err = tx.Table("users").Where(entity.Model{ID: follow.FollowId}).Update("follower_count", gorm.Expr(followerExpr)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (relationRepo *RelationRepo) FollowList(userId int64) (userList []*entity.User, err error) {
	var follow entity.Follow
	err = global.DB.
		Model(follow).
		Select("users.avatar", "users.name").
		Joins("LEFT JOIN users on follows.follow_id = users.id").
		Where(entity.Follow{UserId: userId, IsFollow: true}).
		Find(&userList).Error
	if err != nil {
		return nil, err
	}

	return userList, err
}

func (relationRepo *RelationRepo) FollowerList(userId int64) (userList []*entity.User, err error) {
	var follow entity.Follow
	err = global.DB.Model(follow).
		Select("users.avatar", "users.name").
		Joins("LEFT JOIN users on follows.follow_id = users.id").
		Where(entity.Follow{FollowId: userId, IsFollow: true}).
		Find(&userList).Error
	if err != nil {
		return nil, err
	}
	return userList, err
}

func (relationRepo *RelationRepo) FriendList(userId int64) (userList []*entity.User, err error) {
	err = global.DB.Raw("SELECT users.id,users.`name`,users.avatar FROM "+
		"(SELECT f1.follow_id  FROM ( SELECT * FROM follows WHERE user_id = ? ) AS f1 "+
		"INNER JOIN follows f2 ON f1.follow_id = f2.user_id ) AS friend "+
		"LEFT JOIN users ON friend.follow_id = users.id", userId).Scan(&userList).Error
	if err != nil {
		return nil, err
	}
	return userList, err
}

func (relationRepo *RelationRepo) CheckFollow(userId, targetId int64) (bool, error) {
	var follow *entity.Follow
	var count int64
	err := global.DB.Model(follow).Where(&entity.Follow{UserId: userId, FollowId: targetId, IsFollow: true}).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (relationRepo *RelationRepo) CheckFriend(userId, targetId int64) (bool, error) {
	var cnt int
	err := global.DB.Raw("SELECT COUNT(*) FROM (SELECT * FROM follows WHERE user_id = ?) AS f1 "+
		"INNER JOIN follows f2 ON f1.follow_id = f2.user_id where f1.follow_id = ?", userId, targetId).Scan(&cnt).Error
	if err != nil || cnt == 0 {
		return false, err
	}
	return true, nil
}
