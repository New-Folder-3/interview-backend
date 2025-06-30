package op

import (
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
)

func CreateConversation(UserID string, c *Conversation, preferRole int) (string, error) {
	id := util.GenerateToken(16)
	conversation := model.Conversation{
		ID:         id,
		UserID:     UserID,
		Model:      c.Model,
		PreferRole: preferRole,

		ResultFormat:      c.Parameters.ResultFormat,
		Temperature:       c.Parameters.Temperature,
		TopP:              c.Parameters.TopP,
		EnableThinking:    c.Parameters.EnableThinking,
		PresencePenalty:   c.Parameters.PresencePenalty,
		IncrementalOutput: c.Parameters.IncrementalOutput,
	}
	err := db.CreateConversation(&conversation)
	if err != nil {
		return "", errors.WithStack(err)
	}
	user, err := db.GetUser(UserID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	user.Conversation = util.StringListToDB(append(util.DBToStringList(user.Conversation), id))
	err = db.UpdateUser(user)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return id, nil
}

func CreateMessage(ConversationID string, role int) (string, error) {
	id := util.GenerateToken(16)
	message := model.Message{
		ID:             id,
		Role:           conf.Role[role],
		ConversationID: ConversationID,
	}
	err := db.CreateMessage(&message)
	if err != nil {
		return "", errors.WithStack(err)
	}
	conversation, err := db.GetConversation(ConversationID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	conversation.Messages = util.StringListToDB(append(util.DBToStringList(conversation.Messages), id))
	err = db.UpdateConversation(conversation)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return id, nil
}

func CreateContent(MessageID string, c *Content) (string, error) {
	id := util.GenerateToken(16)
	var videoPtr *string
	if c.Video != nil {
		videos := util.StringListToDB(*c.Video)
		videoPtr = &videos
	} else {
		videoPtr = nil
	}
	content := model.Content{
		ID:        id,
		MessageID: MessageID,
		Text:      c.Text,
		Audio:     c.Audio,
		Image:     c.Image,
		Video:     videoPtr,
	}
	err := db.CreateContent(&content)
	if err != nil {
		return "", errors.WithStack(err)
	}
	message, err := db.GetMessage(MessageID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	message.Contents = util.StringListToDB(append(util.DBToStringList(message.Contents), id))
	err = db.UpdateMessage(message)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return id, nil
}

func AddContent(MessageID, text string, images, audios []string, videos [][]string) error {
	if text != "" {
		_, err := CreateContent(MessageID, &Content{
			Text: &text,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	for _, image := range images {
		_, err := CreateContent(MessageID, &Content{
			Image: &image,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	for _, audio := range audios {
		_, err := CreateContent(MessageID, &Content{
			Audio: &audio,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	for _, video := range videos {
		if len(video) == 0 {
			continue
		}
		_, err := CreateContent(MessageID, &Content{
			Video: &video,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
