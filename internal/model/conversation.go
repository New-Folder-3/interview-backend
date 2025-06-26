package model

type Conversation struct {
	ID         string `json:"id" gorm:"primary_key;unique;not null"`
	UserID     string `json:"userid" form:"userid"`
	Model      string `json:"model" form:"model"`
	Messages   string `json:"talkings" form:"talkings"` // IDs of the talking entries in this conversation
	PreferRole int    `json:"prefer_role" form:"prefer_role"`

	ResultFormat      *string  `json:"result_format,omitempty"`
	Temperature       *float64 `json:"temperature,omitempty"`
	TopP              *float64 `json:"top_p,omitempty"`
	EnableThinking    *bool    `json:"enable_thinking,omitempty"`
	PresencePenalty   *float64 `json:"presence_penalty,omitempty"`
	IncrementalOutput *bool    `json:"incremental_output,omitempty"`
}
