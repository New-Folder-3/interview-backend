package handles

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/conf"
	"interview-backend/internal/db"
	"interview-backend/op"
	"interview-backend/op/chat"
	"interview-backend/util"
)

func GetAllConversations(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	user, _ := c.Get("user")
	if user.(string) != request.Username {
		util.ErrorResp(c, "unmatched user info", 500, false)
		return
	}
	userDB, err := db.GetUser(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	conversations := util.DBToStringList(userDB.Conversation)
	util.SuccessResp(c, conversations, "Get all Conversations Successfully")
}

func GetConversation(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	if request.ConversationID == "" {
		util.ErrorResp(c, "No Conversation ID", 400, false)
		return
	}
	conversation, err := op.GetConversation(request.ConversationID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
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
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	userShould, _ := c.Get("user")
	if request.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	if request.Username == "" || request.PreferRole == nil {
		util.ErrorResp(c, "Invalid User or AI Chosen", 500, false)
		return
	}
	prompt := conf.SysPrompt[*request.PreferRole]
	conversation := chat.NewConversation().AddText(prompt, 0)
	conversationID, err := op.CreateConversation(request.Username, conversation, *request.PreferRole)
	request.ConversationID = conversationID
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	messageID, err := op.CreateMessage(conversationID, 0)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	_, err = op.CreateContent(messageID, &chat.Content{
		Text: &prompt,
	})
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	util.SuccessResp(c, request, "Create Conversation Successfully")
}

func DeleteConversation(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400, false)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	if request.ConversationID == "" {
		util.ErrorResp(c, "No Conversation ID", 400, false)
		return
	}
	if err := op.DeleteConversation(request.ConversationID); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500, false)
		return
	}
	util.SuccessResp(c, nil, "Delete Conversation Successfully")
}
