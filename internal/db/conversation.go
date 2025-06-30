package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func GetConversation(ConversationID string) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := db.Where("id = ?", ConversationID).First(&conversation).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &conversation, nil
}

func CreateConversation(c *model.Conversation) error {
	return errors.WithStack(db.Create(c).Error)
}

func DeleteConversation(ConversationID string) error {
	return errors.WithStack(db.Where("id = ?", ConversationID).Delete(&model.Conversation{}).Error)
}

func UpdateConversation(c *model.Conversation) error {
	return errors.WithStack(db.Model(c).Updates(c).Error)
}

func ReplaceConversation(c *model.Conversation) error {
	return errors.WithStack(db.Model(c).Save(c).Error)
}
