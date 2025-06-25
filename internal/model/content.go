package model

type Content struct {
	ID        string  `json:"id" gorm:"primary_key;unique;not null"`
	MessageID string  `json:"messageid" gorm:"not null"`
	Text      *string `json:"text" gorm:"not null"`
	Video     *string `json:"video" gorm:"not null"`
	Image     *string `json:"image" gorm:"not null"`
	Audio     *string `json:"audio" gorm:"not null"`
}
