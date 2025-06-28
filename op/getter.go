package op

import (
	"fmt"
	"github.com/pkg/errors"
	"interview-backend/internal/db"
	"interview-backend/op/chat"
	"interview-backend/util"
)

func GetConversation(ConversationID string) (*chat.Conversation, error) {
	conversation, err := db.GetConversation(ConversationID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	messageIDs := util.DBToStringList(conversation.Messages)
	messages, err := GetMessage(messageIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret := &chat.Conversation{
		Model: conversation.Model,
		Input: struct {
			Messages []chat.Message `json:"messages"`
		}{Messages: *messages},
		Parameters: &chat.Parameters{
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

func GetMessage(messageIDs []string) (*[]chat.Message, error) {
	var ret []chat.Message
	for _, messageID := range messageIDs {
		message, err := db.GetMessage(messageID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		contentIDs := util.DBToStringList(message.Contents)
		contents, err := GetContent(contentIDs)
		fmt.Println(messageID, contentIDs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ret = append(ret, chat.Message{
			Role:    message.Role,
			Content: *contents,
		})
	}
	return &ret, nil
}

func GetContent(contentIDs []string) (*[]chat.Content, error) {
	var ret []chat.Content
	for _, contentID := range contentIDs {
		content, err := db.GetContent(contentID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var videoPtr *[]string
		if content.Video != nil {
			videos := util.DBToStringList(*content.Video)
			videoPtr = &videos
		} else {
			videoPtr = nil
		}
		ret = append(ret, chat.Content{
			Text:  content.Text,
			Image: content.Image,
			Video: videoPtr,
			Audio: content.Audio,
		})
	}
	return &ret, nil
}
