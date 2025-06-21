package model

type User struct {
	UUID    string `json:"uuid" gorm:"primaryKey"`
	Email   string `json:"email" gorm:"unique"`
	Phone   string `json:"phone" gorm:"unique"`
	PwdHash string `json:"-"`
	PwdTS   int64  `json:"-"`
	// 0: admin
	// 1: user
	Role int `json:"role"`
}
