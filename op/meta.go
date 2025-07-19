package op

type Conversation struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Modalities []string  `json:"modalities"`
	Audio      Audio     `json:"audio,omitempty"`  // 语音合成的声音类型和格式
	Stream     bool      `json:"stream,omitempty"` // 是否开启流式输出
}

type Audio struct {
	Voice  string `json:"voice,omitempty"`  // 语音合成的声音类型
	Format string `json:"format,omitempty"` // 音频格式
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`

	Text       *string     `json:"text,omitempty"`
	ImageURL   *ImageURL   `json:"image_url,omitempty"`
	InputAudio *InputAudio `json:"input_audio,omitempty"`
	VideoURL   *VideoURL   `json:"video_url,omitempty"`
}

type VideoURL struct {
	URL string `json:"url"`
}

type InputAudio struct {
	Format string `json:"format"`
	Data   string `json:"data"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Response struct {
	Choices           []Choice `json:"choices"`
	Object            string   `json:"object"`
	Usage             any      `json:"usage"` // 用 interface{} 也可以
	Created           int64    `json:"created"`
	SystemFingerprint *string  `json:"system_fingerprint"`
	Model             string   `json:"model"`
	ID                string   `json:"id"`
}

type Choice struct {
	FinishReason *string `json:"finish_reason"`
	Logprobs     any     `json:"logprobs"`
	Delta        Delta   `json:"delta"`
	Index        int     `json:"index"`
}

type Delta struct {
	Content string        `json:"content,omitempty"` // 文本内容
	Audio   AudioResponse `json:"audio,omitempty"`   // 音频内容
}

type AudioResponse struct {
	Transcript string `json:"transcript,omitempty"` // 音频转文本结果
	Data       string `json:"data,omitempty"`       // 音频数据
}
