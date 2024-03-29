package repository

import (
	"errors"
	"gorm.io/gorm"
	"tiktok/global"
	"tiktok/model/entity"
)

type CommentRepo struct{}

func (commentRepo *CommentRepo) CommentList(videoId int64) (commentList []*entity.Comment, err error) {
	comment := entity.Comment{VideoId: videoId}
	err = global.DB.Model(comment).Where(comment).Order("create_date DESC").Scan(&commentList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return commentList, nil
}

func (commentRepo *CommentRepo) AddComment(comment *entity.Comment) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		comment.GenerateID()
		err := tx.Create(&comment).Error
		if err != nil {
			return err
		}
		err = tx.Table("videos").Where(&entity.Model{ID: comment.VideoId}).Update("comment_count", gorm.Expr("comment_count + 1")).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (commentRepo *CommentRepo) DelComment(commentId int64) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		comment := entity.Comment{
			Model: entity.Model{ID: commentId},
		}
		err := tx.Model(comment).Delete(&comment).Error
		if err != nil {
			return err
		}
		err = tx.Table("videos").Where(&entity.Model{ID: comment.VideoId}).Update("comment_count", gorm.Expr("comment_count - 1")).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
