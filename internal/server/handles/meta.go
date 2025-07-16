package handles

import "interview/internal/model"

type AuthRequest struct {
	Username  string `json:"username" form:"username"`
	PwdHash   string `json:"pwd" form:"pwd"`
	Token     string `json:"token" form:"token"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
	KeepLogin bool   `json:"keeplogin" form:"keeplogin"`
}

type UserRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	model.UserChangable
	model.UserPwd
	NewPwdHash string `json:"new_pwd_hash" form:"new_pwd_hash"` // 新密码哈希
}

type ConversationRequest struct {
	Username        string   `json:"username" form:"username"`
	ConversationID  string   `json:"conversation_id" form:"conversation_id"`
	ConversationIDs []string `json:"conversation_ids" form:"conversation_ids"`
	MessageID       string   `json:"message_id" form:"message_id"`
	PreferRole      int      `json:"prefer_role" form:"prefer_role"`
	Text            string   `json:"text,omitempty" form:"text"`
	Audio           string   `json:"audio,omitempty" form:"audio"`
	Image           []string `json:"image,omitempty" form:"image"`
	Video           []string `json:"video,omitempty" form:"video"`
}

type UploadRequest struct {
	Base64 string `json:"base64" form:"base64"`
}

type InterviewRequest struct {
	Username       string `json:"username" form:"username"`
	InterviewID    string `json:"interview_id" form:"interview_id"`
	ConversationID string `json:"conversation_id" form:"conversation_id"`
	model.CommentContent
	CommentRole string `json:"comment_role" form:"comment_role"`
	model.InterviewChangeable
}
