package op

import (
	"interview/internal/conf"
	"strings"
)

func NewConversation() *Conversation {
	return &Conversation{
		Model:      "qwen-vl-max-latest",
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
