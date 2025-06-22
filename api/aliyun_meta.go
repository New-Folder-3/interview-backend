package api

type AliyunConversation struct {
	Model string `json:"model"`
	Input struct {
		Messages []interface{} `json:"messages"`
	} `json:"input"`
	Parameters *AliyunParameters `json:"parameters,omitempty"`
}

type AliyunMDMessage struct {
	Role    string          `json:"role"`
	Content []AliyunContent `json:"content"`
}

type AliyunTxtMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AliyunContent struct {
	Text  *string   `json:"text,omitempty"`
	Image *string   `json:"image,omitempty"`
	Video *[]string `json:"video,omitempty"`
	Audio *string   `json:"audio,omitempty"`
}

type AliyunParameters struct {
	ResultFormat      *string  `json:"result_format,omitempty"`
	Temperature       *float64 `json:"temperature,omitempty"`
	TopP              *float64 `json:"top_p,omitempty"`
	EnableThinking    *bool    `json:"enable_thinking,omitempty"`
	PresencePenalty   *float64 `json:"presence_penalty,omitempty"`
	IncrementalOutput *bool    `json:"incremental_output,omitempty"`
}

type AliyunResponse struct {
	StatusCode int          `json:"status_code"`
	RequestID  string       `json:"request_id"`
	Code       string       `json:"code"`
	Message    string       `json:"message"`
	Output     AliyunOutput `json:"output"`
	Usage      Usage        `json:"usage"`
}

type AliyunOutput struct {
	Text         *string        `json:"text"`          // 可能为null，用指针
	FinishReason *string        `json:"finish_reason"` // 可能为null，用指针
	Choices      []AliyunChoice `json:"choices"`
}

type AliyunChoice struct {
	FinishReason string          `json:"finish_reason"`
	Message      AliyunRXMessage `json:"message"`
}

type AliyunRXMessage struct {
	Role    string          `json:"role"`
	Content []AliyunContent `json:"content"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
