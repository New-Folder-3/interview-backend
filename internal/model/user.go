package model

import (
	"fmt"
	"interview/internal/conf"
	"interview/util"
)

type User struct {
	// can't change
	UserNotChangeable
	// can change
	UserChangable
	// can change seperately
	UserPwd
}

type UserNotChangeable struct {
	ID           string `json:"username" gorm:"PrimaryKey;not null;omitempty"` // 用户名
	Interviews   string `json:"interview" gorm:"omitempty"`                    // 面试ID列表，序列化存储
	Conversation string `json:"conversation" gorm:"omitempty"`                 // 所有对话
}

type UserChangable struct {
	Name              string `json:"name" gorm:"not null;omitempty"`         // 姓名
	Email             string `json:"email" gorm:"unique;not null;omitempty"` // 邮箱
	Phone             string `json:"phone" gorm:"unique;not null;omitempty"` // 电话
	PreferInterviewer int    `json:"prefer_interviewer" gorm:"not null;omitempty"`
	Age               int    `json:"age" gorm:"omitempty"`    // 年龄
	Job               string `json:"role" gorm:"omitempty"`   // 职位
	Gender            string `json:"gender" gorm:"omitempty"` // 性别
}

type UserPwd struct {
	PwdHash string `json:"pwd_hash" gorm:"not null;omitempty"` // 密码哈希
	PwdTS   int64  `json:"pwd_ts" gorm:"not null;omitempty"`   // 密码修改时间戳
}

func NewDefaultUser() User {
	return User{
		UserChangable: UserChangable{
			Name: fmt.Sprintf("面试者%s", util.GenerateToken(6)),
			Job:  conf.Conf.Default.Job,
			Age:  conf.Conf.Default.Age,
		},
	}
}
