package op

import (
	"github.com/pkg/errors"
)

func FastTxt(conversationID, txt string) (string, error) {
	conversation, err := GetConversation(conversationID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	response, err := CommonChat(conversation.AddText(txt, 2))
	if len(response.Output.Choices) == 0 {
		err = errors.New("response is empty")
	}
	if err != nil {
		return "", errors.WithStack(err)
	} else {
		return *response.Output.Choices[0].Message.Content[0].Text, nil
	}
}
