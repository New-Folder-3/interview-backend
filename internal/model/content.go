package model

type Content struct {
	ID        string  `json:"id" gorm:"primary_key;unique;not null"`
	MessageID string  `json:"messageid" gorm:"not null"`
	Text      *string `json:"text"`
	Video     *string `json:"video"`
	Image     *string `json:"image"`
	Audio     *string `json:"audio"`
}
