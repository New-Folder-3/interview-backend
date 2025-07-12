package handles

import (
	"github.com/gin-gonic/gin"
	"interview/op"
	"interview/util"
)

func NewMessage(c *gin.Context) {
	var request ConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userShould, _ := c.Get("user")
	if userShould.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	// Add Message to Conversation)
	messageID, err := op.CreateMessage(request.ConversationID, 1)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	// Add Content to Message
	for _, audio := range request.Audio {
		if err = op.CreateContent(messageID, op.Content{Audio: &audio}); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}
	}
	for _, image := range request.Image {
		if err = op.CreateContent(messageID, op.Content{Image: &image}); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}
	}
	for _, video := range request.Video {
		if err = op.CreateContent(messageID, op.Content{Video: &video}); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}
	}
	if err = op.CreateContent(messageID, op.Content{
		Text: &request.Text}); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	// send message to AI
	conversation, err := op.GetConversation(request.ConversationID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	response, err := op.CommonChat(conversation, c)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	// add response to conversation
	messageID, err = op.CreateMessage(request.ConversationID, 2)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	if len(response.Output.Choices) > 0 {
		err = op.CreateContent(messageID, op.Content{
			Text: response.Output.Choices[0].Message.Content[0].Text,
		})
		if err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}
	} else {
		util.ErrorResp(c, "No response from Model", 500)
		return
	}

	// Return the response
	util.SuccessResp(c, response, "Send Message Successfully")
}

func DeleteMessage(c *gin.Context) {
	var message ConversationRequest
	userShould, _ := c.Get("user")
	if err := c.ShouldBindJSON(&message); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	if userShould.(string) != message.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	if message.MessageID == "" {
		util.ErrorResp(c, "No Message ID", 400)
		return
	}
	if err := op.DeleteMessage(message.MessageID, false); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, nil, "Delete Message Successfully")
}
