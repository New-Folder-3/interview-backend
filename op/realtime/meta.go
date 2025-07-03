package realtime

type RealtimeRequest struct {
	EventID string `json:"event_id"`
	Type    string `json:"type"`

	Audio    string          `json:"audio"`
	Image    string          `json:"image"`
	Session  *RequestSession `json:"session"`
	Response *interface{}    `json:"response"`
}

type RequestSession struct {
	Modalities        []string             `json:"modalities"`
	Voice             string               `json:"voice"`
	InputAudioFormat  string               `json:"input_audio_format"`
	OutputAudioFormat string               `json:"output_audio_format"`
	TurnDetection     RequestTurnDetection `json:"turn_detection"`
}

type RequestTurnDetection struct {
	Type              string  `json:"type"`
	Threshold         float64 `json:"threshold"`
	SilenceDurationMS int     `json:"silence_duration_ms"`
	CreateResponse    bool    `json:"create_response"`
	InterruptResponse bool    `json:"interrupt_response"`
}

type RealtimeResponse struct {
	EventID      string `json:"event_id"`
	Type         string `json:"type"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`

	AudioStartMS int    `json:"audio_start_ms"`
	AudioEndMS   int    `json:"audio_end_ms"`
	Transcript   string `json:"transcript"`
	Delta        string `json:"delta"`
	Text         string `json:"text"`

	Session  *interface{} `json:"session,omitempty"`
	Item     *interface{} `json:"item,omitempty"`
	Error    *interface{} `json:"error,omitempty"`
	Response *interface{} `json:"response,omitempty"`
	Part     *Part        `json:"part,omitempty"`
}

type Part struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
