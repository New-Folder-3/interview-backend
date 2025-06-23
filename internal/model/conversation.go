package model

type Conversation struct {
	UUID      string `json:"uuid" form:"uuid"`
	UserID    string `json:"userid" form:"userid"`
	SysPrompt string `json:"sysprompt" form:"sysprompt"`
}

type Message struct {
	Role  string    `json:"role" form:"role"`
	Text  *string   `json:"text,omitempty" form:"text"`
	Image *string   `json:"image,omitempty" form:"image"`
	Audio *string   `json:"audio,omitempty" form:"audio"`
	Video *[]string `json:"video,omitempty" form:"video"`
}
