package handles

type AuthRequest struct {
	Username  string `json:"username" form:"username"`
	PwdHash   string `json:"pwd" form:"pwd"`
	Token     string `json:"token" form:"token"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
	KeepLogin bool   `json:"keeplogin" form:"keeplogin"`
}

type UserRequest struct {
	Username string `json:"username" form:"username"`
	Nickname string `json:"nickname" form:"nickname"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	OldPwd   string `json:"oldpwd" form:"oldpwd"`
	NewPwd   string `json:"newpwd" form:"newpwd"`
}

type MessageRequest struct {
	ConversationID string     `json:"conversation_id" form:"conversation_id,required"`
	Text           string     `json:"text,omitempty" form:"text"`
	Audio          []string   `json:"audio,omitempty" form:"audio"`
	Image          []string   `json:"image,omitempty" form:"image"`
	Video          [][]string `json:"video,omitempty" form:"video"`
}
