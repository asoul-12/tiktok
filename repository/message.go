package repository

import (
	"tiktok/global"
	"tiktok/model/entity"
	"tiktok/model/query"
)

type MessageRepo struct{}

func (messageRepo *MessageRepo) MessageList(messageQuery query.MessageQuery) (messageList []*entity.Message, err error) {
	var message entity.Message
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

func (messageRepo *MessageRepo) SendMessage(message *entity.Message) (err error) {
	message.GenerateID()
	err = global.DB.Model(message).Create(message).Error
	if err != nil {
		return err
	}
	return nil
}

func (messageRepo *MessageRepo) GetLatestMessage(userId, targetId int64) (message *entity.Message, err error) {
	db := global.DB
	whereExpr := " from_user_id = ? and to_user_id = ? "
	selectExpr := " from_user_id , content , create_time "
	//err = db.Raw("? UNION ? ORDER BY create_time LIMIT 1",
	//	db.Select(selectExpr).
	//		Table("messages").
	//		Where(whereExpr, userId, targetId),
	//	db.Select(selectExpr).
	//		Table("messages").
	//		Where(whereExpr, targetId, userId)).Scan(&message).Error
	sqlExpr := " SELECT " + selectExpr + " FROM messages WHERE " + whereExpr
	err = db.Raw(sqlExpr+" UNION "+sqlExpr+" ORDER BY create_time DESC LIMIT 1 ", userId, targetId, targetId, userId).Scan(&message).Error
	if err != nil {
		return nil, err
	}
	return message, nil

}
