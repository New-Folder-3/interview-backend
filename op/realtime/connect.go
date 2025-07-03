package realtime

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"interview/cmd/flags"
	"interview/internal/conf"
	"interview/op"
	"interview/util"
	"net/url"
	"path"
)

var (
	userAudioName = ""
)

func RemoteToFrontend(done chan error, remote *websocket.Conn, frontend *websocket.Conn, conversationIDChan chan string, userID string) {
	defer close(done)
	conversation := op.NewConversation()
	conversationID, err := op.CreateConversation(userID, conversation, 0)
	conversationIDChan <- conversationID
	if err != nil {
		done <- errors.WithStack(err)
		return
	}
	var ret RealtimeResponse
	for {
		err = remote.ReadJSON(&ret)
		if err != nil || !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			done <- errors.WithStack(err)
			return
		}
		switch ret.Type {
		case "input_audio_buffer.speech_started":
			userAudioName = uuid.NewString() + ".wav"
		case "input_audio_buffer.speech_stopped":
			audioName := userAudioName
			userAudioName = ""
			userAudioURL, _ := url.JoinPath(conf.Conf.Schema.URL + "download/audio" + audioName)
			err = util.StreamWritePCMBase64ToWav(path.Join(flags.DataDir, "audio", audioName), "", true)
			if err != nil {
				util.ErrorPrinter(errors.WithStack(err))
			} else {
				conversation.AddAudio(userAudioURL)
			}
		case "response.text.done":
			conversation.AddText(ret.Text, 2)
		case "response.audio_transcript.done":
			conversation.AddText(ret.Part.Text, 2)
		}
		err = frontend.WriteJSON(ret)
		if err != nil || !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			done <- errors.WithStack(err)
			return
		}
	}
}

func FrontendToRemote(done chan error, remote *websocket.Conn, frontend *websocket.Conn) {
	defer close(done)
	for {
		var ret RealtimeRequest
		err := frontend.ReadJSON(&ret)
		if err != nil || !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			done <- errors.WithStack(err)
			return
		}
		switch ret.Type {
		case "input_audio_buffer.append":
			if userAudioName != "" {
				_ = util.StreamWritePCMBase64ToWav(path.Join(flags.DataDir, "audio", userAudioName), ret.Audio, false)
			}
		}
		err = remote.WriteJSON(ret)
		if err != nil || !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			done <- errors.WithStack(err)
			return
		}
	}
}
