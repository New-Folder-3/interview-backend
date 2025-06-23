package handles

type AuthRequest struct {
	Username  string `json:"username" form:"username"`
	PwdHash   string `json:"pwdhash" form:"pwdhash"`
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
	OldPwd   string `json:"pwdhash" form:"oldpwd"`
	NewPwd   string `json:"newpwd" form:"newpwd"`
}
