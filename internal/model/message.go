package model

type Message struct {
	ID             string `json:"id" gorm:"primary_key;unique;not null"`
	ConversationID string `json:"conversationid" gorm:"not null"`
	Role           string `json:"role"`
	Contents       string `json:"contents"`
}
