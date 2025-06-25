package handles

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/conf"
	"interview-backend/internal/db"
	"interview-backend/op"
	"interview-backend/op/aliyun"
	"interview-backend/util"
)

func GetAllConversations(c *gin.Context) {
	user, _ := c.Get("user")
	if user == nil {
		util.ErrorResp(c, "Invalid", 500, false)
		return
	}
	userDB, err := db.GetUser(user.(string))
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	conversations := util.DBToStringList(userDB.Conversation)
	util.SuccessResp(c, conversations, "Get all Conversations Successfully")
}

func GetConversation(c *gin.Context) {
	conversationID, _ := c.Get("id")
	if conversationID == nil {
		util.ErrorResp(c, "Invalid Conversation ID", 400, false)
		return
	}
	conversation, err := op.GetConversation(conversationID.(string))
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	util.SuccessResp(c, conversation, "Get Conversation Successfully")
}

func CreateConversation(c *gin.Context) {
	user, _ := c.Get("user")
	chosen, _ := c.Get("chosen")
	if user == nil || chosen == nil {
		util.ErrorResp(c, "Invalid User or AI Chosen", 500, false)
		return
	}
	conversation := op.NewConversation().AddText(conf.SysPrompt[chosen.(int)], 0)
	conversationID, err := op.CreateConversation(user.(string), conversation)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	response, err := aliyun.ChatWithProxy(conversation, conf.Conf.API.AliyunAPIKey, c)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	MessageID, err := op.CreateMessage(conversationID, 0)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	_, err = op.CreateContent(MessageID, conversation.AddResponse(response))
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	util.SuccessResp(c, conversationID, "Create Conversation Successfully")
}

func CreateMessage(c *gin.Context) {
	var message MessageRequest
	if err := c.ShouldBindJSON(&message); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	MessageID, err := op.CreateMessage(message.ConversationID, 1)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	conversation, err := op.GetConversation(message.ConversationID)
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
	var response *op.Response
	switch {
	case len(message.Audio) > 0:
		response, err = aliyun.ChatWithProxy(conversation.AddAudio(txts...), conf.Conf.API.AliyunAPIKey, c)
	case len(message.Image) > 0:
		response, err = aliyun.ChatWithProxy(conversation.AddImg(txts...), conf.Conf.API.AliyunAPIKey, c)
	default:
		response, err = aliyun.ChatWithProxy(conversation.AddText(message.Text, 1), conf.Conf.API.AliyunAPIKey, c)
	}
	_, err = op.CreateContent(MessageID, conversation.AddResponse(response))
	if err != nil {
		util.ErrorPrinter(err)
	}
}
