package handles

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/op"
	"interview/util"
	"sync"
	"time"
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

	// Add Message to Conversation
	messageID, err := op.CreateMessage(request.ConversationID, 1)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	// Add Content to Message
	transcription := make([]string, len(request.Audio))
	var transcriptWG sync.WaitGroup
	for index, audio := range request.Audio {
		if err = op.CreateContent(messageID, op.Content{Audio: &audio}); err != nil {
			util.ErrorPrinter(err)
			util.ErrorResp(c, err.Error(), 500)
			return
		}
		transcriptWG.Add(1)
		go func() {
			defer transcriptWG.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			transcript, err := util.ChatSTT(ctx, audio)
			if err != nil {
				transcription[index] = ""
				util.ErrorPrinter(err)
			} else {
				transcription[index] = transcript
			}
		}()
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
	response, err := op.CommonChat(conversation)
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

	ttsURL, err := util.ChatTTS(*response.Output.Choices[0].Message.Content[0].Text)
	if err != nil {
		util.ErrorPrinter(err)
	}

	// Wait for all transcripts to finish
	transcriptWG.Wait()

	// Return the response
	util.SuccessResp(c, map[string]interface{}{
		"text":               *response.Output.Choices[0].Message.Content[0].Text,
		"user_transcription": transcription,
		"tts_url":            ttsURL,
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
