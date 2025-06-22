package api

import (
	"interview-backend/internal/conf"
	"strings"
)

func NewAliyunConversation(role int) *AliyunConversation {
	return &AliyunConversation{
		Model: "qwen-vl-max-latest",
		Input: struct {
			Messages []interface{} `json:"messages"`
		}{
			Messages: []interface{}{
				AliyunTxtMessage{
					Role:    "system",
					Content: conf.SysPrompt[role],
				},
			},
		},
	}
}

func (c *AliyunConversation) ChangeModel(model string) *AliyunConversation {
	c.Model = model
	return c
}

func (c *AliyunConversation) AddText(txt string) *AliyunConversation {
	msg := &AliyunTxtMessage{
		Role:    "user",
		Content: txt,
	}
	c.Input.Messages = append(c.Input.Messages, msg)
	return c
}

func (c *AliyunConversation) AddImg(txts ...string) *AliyunConversation {
	msg := AliyunMDMessage{
		Role: "user",
	}
	for _, txt := range txts {
		content := AliyunContent{}
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

func (c *AliyunConversation) AddAudio(txts ...string) *AliyunConversation {
	msg := AliyunMDMessage{
		Role:    "user",
		Content: make([]AliyunContent, len(txts)),
	}
	for _, txt := range txts {
		content := AliyunContent{}
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
