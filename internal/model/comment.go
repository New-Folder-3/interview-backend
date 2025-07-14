package model

type Comment struct {
	ID string `json:"id" gorm:"primary_key;unique;not null"`
	CommentContent
}

type CommentContent struct {
	Good  string `json:"good" gorm:"not null"`
	Bad   string `json:"bad" gorm:"not null"` // 面试者的评价
	Other string `json:"other" gorm:"not null"`
}
