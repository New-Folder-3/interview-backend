package op

type Conversation struct {
	Model string `json:"model"`
	Input struct {
		Messages []Message `json:"messages"`
	} `json:"input"`
	Parameters *Parameters `json:"parameters,omitempty"`
}

type Parameters struct {
	ResultFormat      *string  `json:"result_format,omitempty"`
	Temperature       *float64 `json:"temperature,omitempty"`
	TopP              *float64 `json:"top_p,omitempty"`
	EnableThinking    *bool    `json:"enable_thinking,omitempty"`
	PresencePenalty   *float64 `json:"presence_penalty,omitempty"`
	IncrementalOutput *bool    `json:"incremental_output,omitempty"`
}

type Response struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Output     Output `json:"output"`
	Usage      Usage  `json:"usage"`
}

type Output struct {
	Text         *string  `json:"text"`          // 可能为null，用指针
	FinishReason *string  `json:"finish_reason"` // 可能为null，用指针
	Choices      []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Text  *string   `json:"text,omitempty"`
	Image *string   `json:"image,omitempty"`
	Video *[]string `json:"video,omitempty"`
	Audio *string   `json:"audio,omitempty"`
}
