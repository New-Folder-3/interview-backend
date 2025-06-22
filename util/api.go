package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"interview-backend/api"
	"interview-backend/internal/client"
	"interview-backend/internal/conf"
	"net/http"
	"strings"
)

func AliyunChatWithProxy(sender *api.AliyunConversation, key string, c *gin.Context) (*api.AliyunResponse, error) {
	jsonPayload, _ := json.Marshal(*sender)
	req, _ := http.NewRequest(http.MethodPost, conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("X-DashScope-SSE", "enable")
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data = strings.TrimPrefix(line, "data:")
			c.Writer.Write([]byte(line + "\n"))
			c.Writer.Flush()
			fmt.Println(data)
		}
	}
	c.Writer.Write([]byte("data: [DONE]"))

	var ret api.AliyunResponse
	err = json.Unmarshal([]byte(data), &ret)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &ret, nil
}
