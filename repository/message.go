package repository

import (
	"tiktok/global"
	"tiktok/model"
	"tiktok/model/query"
)

type MessageRepo struct{}

func (messageRepo *MessageRepo) MessageList(messageQuery query.MessageQuery) (messageList []*model.Message, err error) {
	var message model.Message
	selectExpr := "from_user_id , to_user_id , content , create_time"
	whereExpr := "from_user_id = ? AND to_user_id = ? AND create_time > ?"
	db := global.DB
	err = db.Raw("? UNION ?",
		db.Select(selectExpr).
			Model(message).
			Where(whereExpr, messageQuery.FromUserId, messageQuery.ToUserId, messageQuery.CreateTime),
		db.Select(selectExpr).
			Model(message).
			Where(whereExpr, messageQuery.ToUserId, messageQuery.FromUserId, messageQuery.CreateTime)).
		Order("create_time").
		Scan(&messageList).Error
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

func (messageRepo *MessageRepo) SendMessage(message *model.Message) (err error) {
	message.GenerateID()
	err = global.DB.Model(message).Create(message).Error
	if err != nil {
		return err
	}
	return nil
}
