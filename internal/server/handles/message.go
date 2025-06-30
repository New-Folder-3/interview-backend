package handles

import (
	"github.com/gin-gonic/gin"
	"interview/internal/conf"
	"interview/op"
	"interview/util"
)

func CreateMessage(c *gin.Context) {
	var message ConversationRequest
	if err := c.ShouldBindJSON(&message); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	userShould, _ := c.Get("user")
	if userShould.(string) != message.Username {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	MessageID, err := op.CreateMessage(message.ConversationID, 1)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	var txts []string
	txts = append(txts, message.Image...)
	txts = append(txts, message.Audio...)
	txts = append(txts, message.Text)
	err = op.AddContent(MessageID, message.Text, message.Image, message.Audio, message.Video)
	conversation, err := op.GetConversation(message.ConversationID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}

	var response *op.Response
	switch {
	case len(message.Audio) > 0:
		response, err = op.CommonChat(conversation.AddAudio(txts...), conf.Conf.API.AliyunAPIKey, c)
	case len(message.Image) > 0:
		response, err = op.CommonChat(conversation.AddImg(txts...), conf.Conf.API.AliyunAPIKey, c)
	default:
		response, err = op.CommonChat(conversation.AddText(message.Text, 1), conf.Conf.API.AliyunAPIKey, c)
	}
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	MessageID, err = op.CreateMessage(message.ConversationID, 2)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	content := response.Output.Choices[0].Message.Content[0]
	_, err = op.CreateContent(MessageID, &content)
	util.SuccessResp(c, response, "Create Conversation Successfully")
}

func DeleteMessage(c *gin.Context) {
	var message ConversationRequest
	userShould, _ := c.Get("user")
	if err := c.ShouldBindJSON(&message); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	if userShould.(string) != message.Username {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	if message.MessageID == "" {
		util.ErrorResp(c, "No Message ID", 400, false)
		return
	}
	if err := op.DeleteMessage(message.MessageID, false); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	util.SuccessResp(c, nil, "Delete Message Successfully")
}
