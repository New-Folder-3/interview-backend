package db

import (
	"github.com/pkg/errors"
	"interview-backend/internal/model"
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
	return errors.WithStack(db.Delete(&model.Conversation{}, ConversationID).Error)
}

func UpdateConversation(c *model.Conversation) error {
	return errors.WithStack(db.Model(c).Updates(c).Error)
}
