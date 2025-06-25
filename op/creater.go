package op

import (
	"github.com/pkg/errors"
	"interview-backend/internal/conf"
	"interview-backend/internal/db"
	"interview-backend/internal/model"
	"interview-backend/util"
)

func CreateConversation(UserID string, c *Conversation) (string, error) {
	id := util.GenerateToken(16)
	conversation := model.Conversation{
		ID:     id,
		UserID: UserID,
		Model:  c.Model,

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
	videos := util.StringListToDB(*c.Video)
	var videoPtr *string
	if len(videos) > 0 {
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
	_, err := CreateContent(MessageID, &Content{
		Text: &text,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	for _, image := range images {
		_, err = CreateContent(MessageID, &Content{
			Image: &image,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	for _, audio := range audios {
		_, err = CreateContent(MessageID, &Content{
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
		_, err = CreateContent(MessageID, &Content{
			Video: &video,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
