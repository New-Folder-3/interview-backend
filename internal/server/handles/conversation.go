package handles

import (
	"github.com/gin-gonic/gin"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/op"
	"interview/util"
)

func GetAllConversations(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	user, _ := c.Get("user")
	if user.(string) != request.Username {
		util.ErrorResp(c, "unmatched user info", 500)
		return
	}
	userDB, err := db.GetUser(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	conversations := util.DBToStringList(userDB.Conversation)
	util.SuccessResp(c, conversations, "Get all Conversations Successfully")
}

func GetConversation(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	if request.ConversationID == "" {
		util.ErrorResp(c, "No Conversation ID", 400)
		return
	}
	conversation, err := op.GetConversation(request.ConversationID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, conversation, "Get Conversation Successfully")
}

// need username, prefer_role
// provide conversation_id
func CreateConversation(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userShould, _ := c.Get("user")
	if request.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	prompt := conf.SysPrompt[request.PreferRole]
	conversation := op.NewConversation().AddText(prompt, 0)
	conversationID, err := op.CreateConversation(request.Username, conversation, request.PreferRole)
	request.ConversationID = conversationID
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	messageID, err := op.CreateMessage(conversationID, 0)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	_, err = op.CreateContent(messageID, &op.Content{
		Text: &prompt,
	})
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, request, "Create Conversation Successfully")
}

func DeleteConversation(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	if request.ConversationID == "" {
		util.ErrorResp(c, "No Conversation ID", 400)
		return
	}
	if err := op.DeleteConversation(request.ConversationID); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, nil, "Delete Conversation Successfully")
}
