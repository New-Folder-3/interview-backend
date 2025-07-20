package model

import "interview/util"

type Interview struct {
	ID       string `json:"id" gorm:"primary_key;unique;not null"` // 面试ID
	Username string `json:"username" form:"userid"`                // 面试者ID
	CreateAt int64  `json:"create_at"`                             // 开始的时间戳

	InterviewChangeable
	ModelComments      string             `json:"model_comments"`                                                // 模型评价
	UserComments       string             `json:"user_comments"`                                                 // 其他评价
	Keywords           string             `json:"keywords" form:"keywords"`                                      // 关键词，序列
	InterviewDimension InterviewDimension `json:"interview_dimension" gorm:"embedded;embeddedPrefix:dimension_"` // 用户维度评分
}

type InterviewDimension struct {
	Hard        float64 `json:"hard"`        // 硬技能
	Soft        float64 `json:"soft"`        // 软技能
	Potential   float64 `json:"potential"`   // 潜力
	Confidence  float64 `json:"confidence"`  // 自信心
	Development float64 `json:"development"` // 发展潜力
	Fit         float64 `json:"fit"`         // 岗位适配度
}

type InterviewChangeable struct {
	ResumeConversation    string `json:"resume_conversation" form:"resume_conversation"`       // 第一次对话ID
	InterviewConversation string `json:"interview_conversation" form:"interview_conversation"` // 第二次对话ID
	EmotionConversation   string `json:"emotion_conversation" form:"emotion_conversation"`     // 第三次对话ID
	VideoURL              string `json:"video_url" form:"video_url"`
	ResumeURL             string `json:"resume_url" form:"resume_url"`
	VideoClipURL          string `json:"video_url_clip" form:"video_url_clip"` // 视频片段URL
}

func NewDefaultInterview() Interview {
	return Interview{
		ID:       "interview_" + util.GenerateToken(16),
		Username: "",
		InterviewChangeable: InterviewChangeable{
			ResumeConversation:    "",
			InterviewConversation: "",
			EmotionConversation:   "",
			VideoURL:              "",
		},
		ModelComments:      "",
		UserComments:       "",
		Keywords:           "",
		InterviewDimension: InterviewDimension{},
	}
}
