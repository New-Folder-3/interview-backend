package realtime

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/op"
	"interview/util"
	"net/url"
	"path"
)

var userAudioName map[string]string

func RemoteToFrontend(done chan error, remote *websocket.Conn, frontend *websocket.Conn, conversationIDChan chan string, userID string) {
	defer close(done)
	conversation := op.NewConversation(conf.Conf.Model.ChatModel, conf.Conf.Default.ModelVoice)
	conversationID, err := op.CreateConversation(userID, conversation, 0)
	conversationIDChan <- conversationID
	if err != nil {
		done <- errors.WithStack(err)
		return
	}
	var ret Response
	for {
		msgType, ori, err := remote.ReadMessage()
		if err != nil {
			done <- errors.WithStack(err)
			return
		}
		if msgType == websocket.TextMessage {
			err = json.Unmarshal(ori, &ret)
			if err != nil {
				util.ErrorPrinter(errors.WithStack(err))
			}
			switch ret.Type {
			case "input_audio_buffer.speech_started":
				userAudioName[conversationID] = uuid.NewString() + ".wav"
			case "input_audio_buffer.speech_stopped":
				audioName := userAudioName[conversationID]
				userAudioName[conversationID] = ""
				userAudioURL, _ := url.JoinPath(conf.Conf.Schema.URL + "download/audio" + audioName)
				err = util.StreamWritePCMBase64ToWav(path.Join(flags.DataDir, "audio", audioName), "", true)
				if err != nil {
					util.ErrorPrinter(errors.WithStack(err))
				} else {
					conversation.AddAudio(userAudioURL, "")
				}
			case "response.text.done":
				conversation.AddText(ret.Text, 2)
			case "response.audio_transcript.done":
				conversation.AddText(ret.Part.Text, 2)
			}
		}
		err = frontend.WriteMessage(msgType, ori)
		if err != nil {
			done <- errors.WithStack(err)
			return
		}
	}
}

func FrontendToRemote(done chan error, remote *websocket.Conn, frontend *websocket.Conn, conversationID string) {
	defer close(done)
	var ret Request
	for {
		msgType, ori, err := frontend.ReadMessage()
		if err != nil {
			done <- errors.WithStack(err)
			return
		}
		audioName, ok := userAudioName[conversationID]
		if msgType == websocket.TextMessage {
			err = json.Unmarshal(ori, &ret)
			if err != nil {
				util.ErrorPrinter(errors.WithStack(err))
			}
			util.StructPrinter(ret, 0)
			switch ret.Type {
			case "input_audio_buffer.append":
				if audioName != "" && ok {
					_ = util.StreamWritePCMBase64ToWav(path.Join(flags.DataDir, "audio", audioName), ret.Audio, false)
				}
			}
		}
		log.Printf("Sending message type:%d size:%d\n", msgType, len(ori))
		err = remote.WriteMessage(msgType, ori)
		if err != nil {
			done <- errors.WithStack(err)
			return
		}
	}
}
