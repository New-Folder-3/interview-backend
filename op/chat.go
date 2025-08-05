package op

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"interview/internal/client"
	"interview/internal/conf"
	"io"
	"net/http"
	"strings"
)

func CommonChat(sender *Conversation) (string, []string, error) {
	if conf.Conf.API.AliyunAPIKey == "" {
		return "", nil, errors.New("API Key is required")
	}
	key := conf.Conf.API.AliyunAPIKey
	jsonPayload, _ := json.Marshal(*sender)
	req, _ := http.NewRequest(http.MethodPost, conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return "", nil, errors.WithStack(fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body))
	}

	scanner := bufio.NewScanner(resp.Body)
	ret := ""
	var audio []string
	for scanner.Scan() {
		line := scanner.Text()
		dat, _ := strings.CutPrefix(line, "data: ")
		if dat == "" || dat == "[DONE]" {
			continue
		}
		var response Response
		if err := json.Unmarshal([]byte(dat), &response); err != nil {
			continue
		}
		ret += response.Choices[0].Delta.Content
		audio = append(audio, response.Choices[0].Delta.Audio.Data)
	}
	if err := scanner.Err(); err != nil {
		return "", nil, errors.WithStack(err)
	}
	return ret, audio, nil
}

func FastChatTxt(conversationID, txt string) (string, error) {
	fmt.Println(txt)
	conversation, err := GetConversation(conversationID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	response, _, err := CommonChat(conversation.AddText(txt, 1))
	switch {
	case err != nil:
		return "", errors.WithStack(err)
	case response == "":
		return "", errors.New("response is empty")
	}
	return response, nil
}
