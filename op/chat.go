package op

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"interview/internal/client"
	"interview/internal/conf"
	"io"
	"net/http"
)

func CommonChat(sender *Conversation) (*Response, error) {
	if conf.Conf.API.AliyunAPIKey == "" {
		return nil, errors.New("API Key is required")
	}
	key := conf.Conf.API.AliyunAPIKey
	jsonPayload, _ := json.Marshal(*sender)
	req, _ := http.NewRequest(http.MethodPost, conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var ret Response
	err = json.Unmarshal(data, &ret)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case ret.Code != "":
		return nil, errors.WithStack(errors.New(ret.Code))
	case len(ret.Output.Choices) == 0:
		return nil, errors.WithStack(errors.New("No Valid Response"))
	}
	return &ret, nil
}
