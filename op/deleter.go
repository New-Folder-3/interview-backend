package op

import (
	"github.com/pkg/errors"
	"interview/internal/db"
	"interview/util"
)

func DeleteConversation(ConversationID string) error {
	conversation, err := db.GetConversation(ConversationID)
	if err != nil {
		return errors.WithStack(err)
	}
	user, err := db.GetUser(conversation.UserID)
	if err != nil {
		return errors.WithStack(err)
	}
	user.Conversation = util.RemoveFromDBList(user.Conversation, ConversationID)
	if err := db.ReplaceUser(user); err != nil {
		return errors.WithStack(err)
	}

	for _, messageID := range util.DBToStringList(conversation.Messages) {
		if err := DeleteMessage(messageID, true); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := db.DeleteConversation(ConversationID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func DeleteMessage(MessageID string, internal bool) error {
	message, err := db.GetMessage(MessageID)
	if err != nil {
		return errors.WithStack(err)
	}
	if !internal { // remove in father conversation
		conversation, err := db.GetConversation(message.ConversationID)
		if err != nil {
			return errors.WithStack(err)
		}
		conversation.Messages = util.RemoveFromDBList(conversation.Messages, MessageID)
		if err := db.ReplaceConversation(conversation); err != nil {
			return errors.WithStack(err)
		}
	}

	for _, contentID := range util.DBToStringList(message.Contents) {
		if err := DeleteContent(contentID, true); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := db.DeleteMessage(MessageID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func DeleteContent(ContentID string, internal bool) error {
	if !internal { // remove in father message
		content, err := db.GetContent(ContentID)
		if err != nil {
			return errors.WithStack(err)
		}
		message, err := db.GetMessage(content.MessageID)
		if err != nil {
			return errors.WithStack(err)
		}
		message.Contents = util.RemoveFromDBList(message.Contents, ContentID)
		if err := db.ReplaceMessage(message); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := db.DeleteContent(ContentID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
