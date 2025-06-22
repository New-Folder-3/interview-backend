package handles

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	PwdHash  string `json:"pwdhash"`
	Token    string `json:"token"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	PwdHash  string `json:"pwdhash" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}
