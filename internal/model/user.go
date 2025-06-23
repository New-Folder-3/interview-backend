package model

type User struct {
	Name              string `json:"name" gorm:"not null"`
	ID                string `json:"id" gorm:"PrimaryKey;not null"`
	Email             string `json:"email" gorm:"unique;not null"`
	Phone             string `json:"phone" gorm:"unique;not null"`
	PwdHash           string `json:"-" gorm:"not null"`
	PwdTS             int64  `json:"-" gorm:"not null"`
	PreferInterviewer int    `json:"prefer_interviewer" gorm:"not null"`
	// 0: ssss
	// 1: ssss

	Role int `json:"role"`
	// 0: student
	// 1: public
	// 2: programmer
}
