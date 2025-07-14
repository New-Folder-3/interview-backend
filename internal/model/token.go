package model

type Token struct {
	UserID   string `json:"username" gorm:"not null"`
	Token    string `json:"token" gorm:"PrimaryKey"`
	ExpireTS int64  `json:"expire_ts"`
}
