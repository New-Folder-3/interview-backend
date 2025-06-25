package op

import (
	"github.com/pkg/errors"
	"interview-backend/internal/db"
	"interview-backend/util"
)

func GetConversation(ConversationID string) (*Conversation, error) {
	conversation, err := db.GetConversation(ConversationID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	messageIDs := util.DBToStringList(conversation.Messages)
	messages, err := GetMessage(messageIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret := &Conversation{
		Model: conversation.Model,
		Input: struct {
			Messages []Message `json:"messages"`
		}{Messages: *messages},
		Parameters: &Parameters{
			ResultFormat:      conversation.ResultFormat,
			Temperature:       conversation.Temperature,
			TopP:              conversation.TopP,
			EnableThinking:    conversation.EnableThinking,
			PresencePenalty:   conversation.PresencePenalty,
			IncrementalOutput: conversation.IncrementalOutput,
		},
	}
	return ret, nil
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

func GetContent(contentIDs []string) (*[]Content, error) {
	var ret []Content
	for _, contentID := range contentIDs {
		content, err := db.GetContent(contentID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		videos := util.DBToStringList(*content.Video)
		var videoPtr *[]string
		if len(videos) > 0 {
			videoPtr = &videos
		} else {
			videoPtr = nil
		}
		ret = append(ret, Content{
			Text:  content.Text,
			Image: content.Image,
			Video: videoPtr,
			Audio: content.Audio,
		})
	}
	return &ret, nil
}
