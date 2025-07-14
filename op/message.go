package op

import (
	"context"
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
	"sync"
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

func NewMessage(conversationID, Text string, Audio, Image []string, Video [][]string) (string, []string, string, error) {
	// Add Message to Conversation
	messageID, err := CreateMessage(conversationID, 1)
	if err != nil {
		return "", nil, "", errors.WithStack(err)
	}

	// Add Content to Message
	transcription := make([]string, len(Audio))
	var transcriptWG sync.WaitGroup
	for index, audio := range Audio {
		if err = CreateContent(messageID, Content{Audio: &audio}); err != nil {
			return "", nil, "", errors.WithStack(err)
		}
		transcriptWG.Add(1)
		go func() {
			defer transcriptWG.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			transcript, err := util.ChatSTT(ctx, audio)
			if err != nil {
				transcription[index] = ""
				util.ErrorPrinter(err)
			} else {
				transcription[index] = transcript
			}
		}()
	}
	for _, image := range Image {
		if err = CreateContent(messageID, Content{Image: &image}); err != nil {
			return "", nil, "", errors.WithStack(err)
		}
	}
	for _, video := range Video {
		if err = CreateContent(messageID, Content{Video: &video}); err != nil {
			return "", nil, "", errors.WithStack(err)
		}
	}
	if err = CreateContent(messageID, Content{
		Text: &Text}); err != nil {
		return "", nil, "", errors.WithStack(err)
	}

	// send message to AI
	conversation, err := GetConversation(conversationID)
	if err != nil {
		return "", nil, "", errors.WithStack(err)
	}
	response, err := CommonChat(conversation)
	if err != nil {
		return "", nil, "", errors.WithStack(err)
	}

	// add response to conversation
	messageID, err = CreateMessage(conversationID, 2)
	if err != nil {
		return "", nil, "", errors.WithStack(err)
	}
	if len(response.Output.Choices) > 0 {
		err = CreateContent(messageID, Content{
			Text: response.Output.Choices[0].Message.Content[0].Text,
		})
		if err != nil {
			return "", nil, "", errors.WithStack(err)
		}
	} else {
		return "", nil, "", errors.New("response is empty")
	}

	ttsURL, err := util.ChatTTS(*response.Output.Choices[0].Message.Content[0].Text)
	if err != nil {
		util.ErrorPrinter(err)
	}

	// Wait for all transcripts to finish
	transcriptWG.Wait()

	return *response.Output.Choices[0].Message.Content[0].Text, transcription, ttsURL, nil
}
