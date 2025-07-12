package handles

import (
	"fmt"
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
	conversationID, err := op.CreateConversation(request.Username, op.NewConversation(), request.PreferRole)
	request.ConversationID = conversationID
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

func GetDimension(c *gin.Context) {
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

	userDB, err := db.GetUser(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	response, err := op.FastTxt(request.ConversationID, fmt.Sprintf(conf.PromptTemplate[conf.DimensionPrompt], userDB.Job))
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	dimensions := util.ChatStringToFloatSlice(response)
	for index, dimension := range dimensions {
		switch index {
		case 0:
			userDB.UserDimension.Hard = dimension
		case 1:
			userDB.UserDimension.Soft = dimension
		case 2:
			userDB.UserDimension.Potential = dimension
		case 3:
			userDB.UserDimension.Confidence = dimension
		case 4:
			userDB.UserDimension.Development = dimension
		case 5:
			userDB.UserDimension.Fit = dimension
		}
	}

	if err = db.UpdateUser(userDB); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, userDB.UserDimension, "Get Dimension Successfully")
}

func GetKeywords(c *gin.Context) {
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

	userDB, err := db.GetUser(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	response, err := op.FastTxt(request.ConversationID,
		fmt.Sprintf(conf.PromptTemplate[conf.KeywordsPrompt], conf.Keywords))
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	keywords := util.ChatStringToStringSlice(response)
	userDB.Keywords = util.StringListToDB(keywords)

	if err = db.UpdateUser(userDB); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, keywords, "Get Keywords Successfully")
}

func CombineConversations(c *gin.Context) {
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

	var newConversationID string
	var newConversationMessage []string
	firstConversation := true
	for _, conversationID := range request.ConversationIDs {
		conversationDB, err := db.GetConversation(conversationID)
		if err != nil {
			continue
		}

		if firstConversation {
			newConversationID, err = op.CreateConversation(request.Username, op.NewConversation(), conversationDB.PreferRole)
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
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	newConversationDB.Messages = util.StringListToDB(newConversationMessage)
	if err = db.UpdateConversation(newConversationDB); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, newConversationID, "Combine Conversations Successfully")
}
