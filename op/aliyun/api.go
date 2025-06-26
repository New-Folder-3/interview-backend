package aliyun

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"interview-backend/internal/client"
	"interview-backend/internal/conf"
	"interview-backend/op"
	"net/http"
	"strings"
)

func ChatWithProxy(sender *op.Conversation, key string, c *gin.Context) (*op.Response, error) {
	jsonPayload, _ := json.Marshal(*sender)
	req, _ := http.NewRequest(http.MethodPost, conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("X-DashScope-SSE", "enable")
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	scanner := bufio.NewScanner(resp.Body)
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data = strings.TrimPrefix(line, "data:")
			c.Writer.Write([]byte(line + "\n"))
			c.Writer.Flush()
		}
	}
	c.Writer.Write([]byte("data: [DONE]"))

	var ret op.Response
	err = json.Unmarshal([]byte(data), &ret)
	if err != nil || ret.Code != "" {
		return nil, errors.WithMessage(err, ret.Code)
	}
	return &ret, nil
}
