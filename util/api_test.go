package util

import (
	"fmt"
	"interview-backend/api"
	"interview-backend/internal/client"
	"testing"
)

func TestAliyunChat(t *testing.T) {
	client.Init()
	conversation := api.NewAliyunConversation(0).
		AddText("你好").
		AddAudio("https://dashscope.oss-cn-beijing.aliyuncs.com/audios/welcome.mp3", "描述一下这段音频").
		ChangeModel("qwen2-audio-instruct")
	ret, _ := AliyunChatWithProxy(conversation, "sk-f689caeb60204d9aaf99c4569556b4c8", nil)
	fmt.Println(*ret.Output.Choices[0].Message.Content[0].Text)
}
