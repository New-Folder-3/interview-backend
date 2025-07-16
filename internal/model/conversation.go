package model

type Conversation struct {
	ID         string `json:"id" gorm:"primary_key;unique;not null"`
	UserID     string `json:"username" form:"userid"`
	Model      string `json:"model" form:"model"`
	Messages   string `json:"talkings" form:"talkings"` // IDs of the talking entries in this conversation
	PreferRole int    `json:"prefer_role" form:"prefer_role"`
	ModelVoice string `json:"model_voice" form:"model_voice"`
}
