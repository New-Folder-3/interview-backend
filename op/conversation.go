package op

import (
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
	"strings"
)

func NewConversation() *Conversation {
	return &Conversation{
		Model:      conf.Conf.Model.ChatModel,
		Parameters: &Parameters{},
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
				Text: &txt,
			},
		},
	}
	c.Input.Messages = append(c.Input.Messages, msg)
	return c
}

func (c *Conversation) AddImg(txts ...string) *Conversation {
	msg := Message{
		Role: "user",
	}
	for _, txt := range txts {
		content := Content{}
		if strings.HasPrefix(txt, "https://") {
			content.Image = &txt
		} else {
			content.Text = &txt
		}
		msg.Content = append(msg.Content, content)
	}
	c.Input.Messages = append(c.Input.Messages, msg)
	return c
}

func (c *Conversation) AddAudio(txts ...string) *Conversation {
	msg := Message{
		Role:    "user",
		Content: make([]Content, len(txts)),
	}
	for _, txt := range txts {
		content := Content{}
		if strings.HasPrefix(txt, "https://") {
			content.Audio = &txt
		} else {
			content.Text = &txt
		}
		msg.Content = append(msg.Content, content)
	}
	c.Input.Messages = append(c.Input.Messages, msg)
	return c
}

func (c *Conversation) AddResponse(r *Response) *Content {
	content := r.Output.Choices[0].Message.Content[0]
	msg := Message{
		Role:    "assistant",
		Content: []Content{content},
	}
	c.Input.Messages = append(c.Input.Messages, msg)
	return &content
}

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

	messageID, err := CreateMessage(id, 0)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err = CreateContent(messageID, Content{
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
			newConversationID, err = CreateConversation(UserID, NewConversation(), conversationDB.PreferRole)
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
