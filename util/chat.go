package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"interview/internal/client"
	"interview/internal/conf"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TTSRequest struct {
	Model string `json:"model"`
	Input struct {
		Text  string `json:"text"`
		Voice string `json:"voice"`
	} `json:"input"`
}

type TTSResponse struct {
	Output struct {
		Audio struct {
			Url string `json:"url"`
		} `json:"audio"`
	} `json:"output"`
}

type STTRequest struct {
	Model string `json:"model"`
	Input struct {
		FileURLs []string `json:"file_urls"`
	} `json:"input"`
}

type STTSubmitResponse struct {
	Output struct {
		TaskID string `json:"task_id"`
	} `json:"output"`
}

type STTTaskResponse struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		Result     []struct {
			TranscriptionURL string `json:"transcription_url"`
			SubtaskStatus    string `json:"subtask_status"`
		} `json:"result"`
	} `json:"output"`
}

type STTJsonResponse struct {
	Transcripts []struct {
		Text string `json:"text"`
	} `json:"transcripts"`
}

func ChatStringToFloatSlice(str string) []float64 {
	if str == "" {
		return nil
	}
	var result []float64
	items := strings.Split(str, " ")
	for _, item := range items {
		f, err := strconv.ParseFloat(item, 64)
		if err == nil {
			result = append(result, f)
		}
	}
	return result
}

func ChatStringToStringSlice(str string) []string {
	if str == "" {
		return nil
	}
	var result []string
	items := strings.Split(str, " ")
	for _, item := range items {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func ChatTTS(txt string) (string, error) {
	if conf.Conf.API.AliyunAPIKey == "" {
		return "", errors.New("API Key is required")
	}
	key := conf.Conf.API.AliyunAPIKey
	request := TTSRequest{
		Model: conf.Conf.Model.TTSModel,
		Input: struct {
			Text  string `json:"text"`
			Voice string `json:"voice"`
		}{Text: txt, Voice: conf.Conf.Model.TTSVoice},
	}
	jsonPayload, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var ret TTSResponse
	if err = json.Unmarshal(data, &ret); err != nil {
		return "", errors.WithStack(err)
	} else {
		return ret.Output.Audio.Url, nil
	}
}

func ChatSTT(ctx context.Context, audioURL string) (string, error) {
	if conf.Conf.API.AliyunAPIKey == "" {
		return "", errors.New("API Key is required")
	}
	key := conf.Conf.API.AliyunAPIKey
	request := STTRequest{
		Model: conf.Conf.Model.STTModel,
		Input: struct {
			FileURLs []string `json:"file_urls"`
		}{FileURLs: []string{audioURL}},
	}
	jsonPayload, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", conf.AliMDUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var ret STTSubmitResponse
	if err = json.Unmarshal(data, &ret); err != nil {
		return "", errors.WithStack(err)
	}
	taskID := ret.Output.TaskID
	var url string
	for {
		select {
		case <-ctx.Done():
			return "", errors.WithStack(ctx.Err())
		default:
		}
		var status string
		status, url, err = ChatSTTTask(taskID)
		if err != nil {
			ErrorPrinter(err)
			continue
		}
		if status == "SUCCEEDED" {
			break
		} else if status == "FAILED" {
			return "", errors.New("STT task failed")
		}
		time.Sleep(1 * time.Second)
	}
	transcription, err := ChatSTTJson(url)
	if err != nil {
		return "", errors.WithStack(err)
	} else {
		return transcription, nil
	}
}

func ChatSTTTask(taskID string) (string, string, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(conf.AliTaskUrlTemplate, taskID), nil)
	req.Header.Set("Authorization", "Bearer "+conf.Conf.API.AliyunAPIKey)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var ret STTTaskResponse
	if err = json.Unmarshal(data, &ret); err != nil {
		return "", "", errors.WithStack(err)
	} else {
		retURL := ""
		if ret.Output.TaskStatus == "SUCCEEDED" {
			retURL = ret.Output.Result[0].TranscriptionURL
		}
		return ret.Output.TaskStatus, retURL, nil
	}
}

func ChatSTTJson(url string) (string, error) {
	if url == "" {
		return "", errors.New("URL is required")
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var ret STTJsonResponse
	if err = json.Unmarshal(data, &ret); err != nil {
		return "", errors.WithStack(err)
	} else if len(ret.Transcripts) > 0 {
		return ret.Transcripts[0].Text, nil
	}
	return "", errors.New("No transcription found")
}
