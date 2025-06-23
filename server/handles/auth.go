package handles

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/db"
	"interview-backend/internal/model"
	"interview-backend/util"
	"time"
)

func LoginPwdHandle(c *gin.Context) {
	var loginRequest AuthRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	var userInfo string
	switch {
	case len(loginRequest.Username) > 0:
		userInfo = loginRequest.Username
	case len(loginRequest.Email) > 0:
		userInfo = loginRequest.Email
	case len(loginRequest.Phone) > 0:
		userInfo = loginRequest.Phone
	}
	user, err := db.GetUser(userInfo)
	if err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
	}
	if user.PwdHash == loginRequest.PwdHash {
		newToken := model.Token{
			Token:  util.GenerateToken(64),
			UserID: userInfo,
		}
		if loginRequest.KeepLogin {
			newToken.ExpireTS = time.Now().AddDate(0, 1, 0).Unix()
		} else {
			newToken.ExpireTS = time.Now().AddDate(0, 0, 1).Unix()
		}
		if err := db.CreateToken(&newToken); err != nil {
			util.ErrorResp(c, err.Error(), 400, false)
			return
		}
		util.SuccessResp(c, newToken, "Login successful")
	} else {
		util.ErrorResp(c, "Invalid credentials", 400, false)
		return
	}
}

func LoginTokenHandle(c *gin.Context) {
	var loginRequest AuthRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	token, err := db.GetToken(loginRequest.Token)
	if err != nil {
		util.ErrorResp(c, "Invalid credentials", 400, false)
	} else {
		util.SuccessResp(c, token, "Login successful")
	}
}

func LogoutHandle(c *gin.Context) {
	token := c.Query("token")
	user := c.Query("user")
	if err := db.DeleteToken(user, token); err != nil {
		util.ErrorResp(c, "Unable to logout", 400, false)
	} else {
		util.SuccessResp(c, nil, "Logged out successfully")
	}
}

func RegisterHandle(c *gin.Context) {
	var registerRequest AuthRequest
	if err := c.ShouldBind(&registerRequest); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	user := model.User{
		ID:      registerRequest.Username,
		Name:    registerRequest.Username,
		PwdHash: registerRequest.PwdHash,
		Email:   registerRequest.Email,
		Phone:   registerRequest.Phone,
		PwdTS:   time.Now().Unix(),
	}
	if err := db.CreateUser(&user); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "Register successful")
	}
}
