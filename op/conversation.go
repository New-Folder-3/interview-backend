package op

import (
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
	"strings"
)

func NewConversation(model, voice string) *Conversation {
	return &Conversation{
		Model:      model,
		Modalities: []string{"text", "audio"},
		Audio: Audio{
			Voice:  voice,
			Format: "wav",
		},
		Stream: true,
	}
}

func (c *Conversation) ChangeModel(model string) *Conversation {
	c.Model = model
	return c
}

func (c *Conversation) AddText(txt string, role int) *Conversation {
	msg := Message{
		Role: conf.Role[role],
		Content: []Content{
			Content{
				Type: "text",
				Text: &txt,
			},
		},
	}
	c.Messages = append(c.Messages, msg)
	return c
}

func (c *Conversation) AddImg(txts ...string) *Conversation {
	msg := Message{
		Role: "user",
	}
	for _, txt := range txts {
		content := Content{}
		if strings.HasPrefix(txt, "https://") {
			content.Type = "image_url"
			content.ImageURL = &ImageURL{URL: txt}
		} else {
			content.Type = "text"
			content.Text = &txt
		}
		msg.Content = append(msg.Content, content)
	}
	c.Messages = append(c.Messages, msg)
	return c
}

func (c *Conversation) AddAudio(audio string, txt string) *Conversation {
	msg := Message{
		Role: "user",
	}
	audioContent := Content{
		Type: "input_audio",
		InputAudio: &InputAudio{
			Format: "mp3",
			Data:   audio,
		},
	}
	if txt == "" {
		txtContent := Content{
			Type: "text",
			Text: &txt,
		}
		msg.Content = append(msg.Content, audioContent, txtContent)
	}
	c.Messages = append(c.Messages, msg)
	return c
}

func (c *Conversation) AddResponse(text string) *Conversation {
	msg := Message{
		Role: "assistant",
		Content: []Content{
			Content{
				Type: "text",
				Text: &text,
			},
		},
	}
	c.Messages = append(c.Messages, msg)
	return c
}

func CreateConversation(UserID string, c *Conversation, preferRole int) (string, error) {
	id := util.GenerateToken(16)
	conversation := model.Conversation{
		ID:         id,
		UserID:     UserID,
		Model:      c.Model,
		PreferRole: preferRole,
		ModelVoice: c.Audio.Voice,
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

	messageID, err := CreateMessage(id, 0)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err = CreateContent(messageID, Content{
		Type: "text",
		Text: &conf.InterviewerPrompt[preferRole],
	}); err != nil {
		return "", errors.WithStack(err)
	}

	return id, nil
}

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
	ret := NewConversation(conversation.Model, conversation.ModelVoice)
	ret.Messages = *messages
	return ret, nil
}

func CombineConversations(ConversationIDs []string, UserID string) (string, error) {
	var newConversationID string
	var newConversationMessage []string
	firstConversation := true
	for _, conversationID := range ConversationIDs {
		conversationDB, err := db.GetConversation(conversationID)
		if err != nil {
			continue
		}

		if firstConversation {
			newConversationID, err = CreateConversation(UserID,
				NewConversation(conversationDB.Model, conversationDB.ModelVoice),
				conversationDB.PreferRole)
			if err != nil {
				continue
			}
			newConversationDB, err := db.GetConversation(newConversationID)
			if err != nil {
				continue
			}
			newConversationMessage = util.DBToStringList(newConversationDB.Messages)
			firstConversation = false
		}

		conversationMessage := util.DBToStringList(conversationDB.Messages)
		newConversationMessage = append(newConversationMessage, conversationMessage[1:]...)
	}

	newConversationDB, err := db.GetConversation(newConversationID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	newConversationDB.Messages = util.StringListToDB(newConversationMessage)
	if err = db.UpdateConversation(newConversationDB); err != nil {
		return "", errors.WithStack(err)
	}

	return newConversationID, nil
}
