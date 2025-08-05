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
	ID           string `json:"username" gorm:"PrimaryKey;not null"` // 用户名
	Interviews   string `json:"interview"`                           // 面试ID列表，序列化存储
	Conversation string `json:"conversation"`                        // 所有对话

	OTP        string `json:"otp"`
	OTPTS      int64  `json:"otp_ts"` // OTP时间戳
	LoginCount int64  `json:"login_count"`
}

type UserChangable struct {
	Name              string `json:"name" gorm:"not null"`         // 姓名
	Email             string `json:"email" gorm:"unique;not null"` // 邮箱
	Phone             string `json:"phone" gorm:"unique;not null"` // 电话
	PreferInterviewer int    `json:"prefer_interviewer" gorm:"not null"`
	Age               int    `json:"age"`                                             // 年龄
	Job               string `json:"role"`                                            // 职位
	Gender            string `json:"gender"`                                          // 性别
	EnvironmentAudio  bool   `json:"environment_audio" gorm:"not null;default:false"` // 是否开启环境音
	ResumeURL         string `json:"resume_url" form:"resume_url"`
}

type UserPwd struct {
	PwdHash string `json:"pwd_hash" gorm:"not null"` // 密码哈希
	PwdTS   int64  `json:"pwd_ts" gorm:"not null"`   // 密码修改时间戳
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
