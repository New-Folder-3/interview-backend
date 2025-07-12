package model

import (
	"fmt"
	"interview/internal/conf"
	"interview/util"
)

type User struct {
	Name              string `json:"name" gorm:"not null"`          // 姓名
	ID                string `json:"id" gorm:"PrimaryKey;not null"` // 用户名
	Email             string `json:"email" gorm:"unique;not null"`  // 邮箱
	Phone             string `json:"phone" gorm:"unique;not null"`  // 电话
	PwdHash           string `json:"-" gorm:"not null"`             // 密码哈希
	PwdTS             int64  `json:"-" gorm:"not null"`             // 密码修改时间戳
	PreferInterviewer int    `json:"prefer_interviewer" gorm:"not null"`
	Conversation      string `json:"conversation"` // 所有对话
	Age               int    `json:"age"`          // 年龄
	Job               string `json:"role"`         // 职位

	Resume    string `json:"resume"`    // 简历对话ID
	Interview string `json:"interview"` // 面试对话ID
	Emotion   string `json:"emotion"`   // 情感对话ID

	Keywords      string        `json:"keywords"`       // 关键词，序列
	UserDimension UserDimension `json:"user_dimension"` // 用户维度评分
}

type UserDimension struct {
	Hard        float64 `json:"hard"`        // 硬技能
	Soft        float64 `json:"soft"`        // 软技能
	Potential   float64 `json:"potential"`   // 潜力
	Confidence  float64 `json:"confidence"`  // 自信心
	Development float64 `json:"development"` // 发展潜力
	Fit         float64 `json:"fit"`         // 岗位适配度
}

func NewDefaultUser() User {
	return User{
		Name: fmt.Sprintf("面试者%s", util.GenerateToken(6)),
		ID:   util.GenerateToken(16),
		Job:  conf.DefaultJob,
		Age:  conf.DefaultAge,
	}
}
