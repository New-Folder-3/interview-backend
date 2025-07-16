package op

import (
	"context"
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
	"time"
)

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

func GetMessage(messageIDs []string) (*[]Message, error) {
	var ret []Message
	for _, messageID := range messageIDs {
		message, err := db.GetMessage(messageID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		contentIDs := util.DBToStringList(message.Contents)
		contents, err := GetContent(contentIDs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ret = append(ret, Message{
			Role:    message.Role,
			Content: *contents,
		})
	}
	return &ret, nil
}

func NewMessage(conversationID, Text string, Audio string, Image []string, Video []string) (string, string, string, error) {
	// Add Message to Conversation
	messageID, err := CreateMessage(conversationID, 1)
	if err != nil {
		return "", "", "", errors.WithStack(err)
	}

	// Add Content to Message
	var transcript string
	transcriptChannel := make(chan struct{})
	if Audio != "" {
		content := Content{
			Type: "input_audio",
			InputAudio: &InputAudio{
				Format: "mp3",
				Data:   Audio,
			},
		}
		if err = CreateContent(messageID, content); err != nil {
			return "", "", "", errors.WithStack(err)
		}

		go func() {
			defer close(transcriptChannel)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			tr, err := util.ChatSTT(ctx, Audio)
			if err != nil {
				util.ErrorPrinter(err)
				transcript = ""
			} else {
				transcript = tr
			}
		}()
	} else {
		close(transcriptChannel)
	}

	for _, image := range Image {
		content := Content{
			Type: "image_url",
			ImageURL: &ImageURL{
				URL: image,
			},
		}
		if err = CreateContent(messageID, content); err != nil {
			return "", "", "", errors.WithStack(err)
		}
	}

	if len(Video) > 0 {
		content := Content{
			Type:  "video",
			Video: &Video,
		}
		if err = CreateContent(messageID, content); err != nil {
			return "", "", "", errors.WithStack(err)
		}
	}

	if Text != "" {
		content := Content{
			Type: "text",
			Text: &Text,
		}
		if err = CreateContent(messageID, content); err != nil {
			return "", "", "", errors.WithStack(err)
		}
	}

	// send message to AI
	conversation, err := GetConversation(conversationID)
	if err != nil {
		return "", "", "", errors.WithStack(err)
	}

	response, rawAudio, err := CommonChat(conversation)
	switch {
	case err != nil:
		return "", "", "", errors.WithStack(err)
	case response == "":
		return "", "", "", errors.New("response is empty")
	}

	// add response to conversation
	messageID, err = CreateMessage(conversationID, 2)
	if err != nil {
		return "", "", "", errors.WithStack(err)
	}

	if err = CreateContent(messageID, Content{
		Type: "text",
		Text: &response,
	}); err != nil {
		return "", "", "", errors.WithStack(err)
	}

	<-transcriptChannel
	audio, err := util.ChunksToWavBase64(rawAudio)
	if err != nil {
		util.ErrorPrinter(err)
		audio = ""
	}
	return response, transcript, audio, nil
}
