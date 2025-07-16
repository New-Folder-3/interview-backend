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

	response, transcription, ttsURL, err := op.NewMessage(request.ConversationID, request.Text, request.Audio, request.Image, request.Video)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	// Return the response
	util.SuccessResp(c, map[string]interface{}{
		"text":               response,
		"user_transcription": transcription,
		"model_audio":        ttsURL,
	}, "Send Message Successfully")
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
