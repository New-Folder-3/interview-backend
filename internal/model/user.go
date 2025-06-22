package model

type User struct {
	Name              string `json:"name" gorm:"not null"`
	ID                string `json:"id" gorm:"PrimaryKey;not null"`
	Email             string `json:"email" gorm:"unique"`
	Phone             string `json:"phone" gorm:"unique"`
	PwdHash           string `json:"-"`
	PwdTS             int64  `json:"-"`
	PreferInterviewer int    `json:"prefer_interviewer" gorm:"not null"`
	// 0: ssss
	// 1: ssss

	Role int `json:"role"`
	// 0: student
	// 1: public
	// 2: programmer
}
