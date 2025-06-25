package handles

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/db"
	"interview-backend/util"
	"time"
)

func GetUser(c *gin.Context) {
	var getRequest UserRequest
	if err := c.ShouldBindJSON(&getRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	userShould, _ := c.Get("user")
	if getRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	user, err := db.GetUser(getRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "User retrieved successfully")
	}
}

func UpdateUserInfo(c *gin.Context) {
	var updateRequest UserRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	userShould, _ := c.Get("user")
	if updateRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	user, err := db.GetUser(updateRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}

	user.Email = updateRequest.Email
	user.Phone = updateRequest.Phone
	user.Name = updateRequest.Nickname

	if err := db.UpdateUser(user); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "User updated successfully")
	}
}

func UpdateUserPwd(c *gin.Context) {
	var updatePwdRequest UserRequest
	if err := c.ShouldBindJSON(&updatePwdRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	userShould, _ := c.Get("user")
	if updatePwdRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400, false)
		return
	}
	
	user, err := db.GetUser(updatePwdRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}

	if updatePwdRequest.OldPwd == user.PwdHash {
		user.PwdHash = updatePwdRequest.NewPwd
		user.PwdTS = time.Now().Unix()
	} else if updatePwdRequest.OldPwd != "" {
		util.ErrorResp(c, "Incorrect Old Password, Password changed failed", 400, false)
		return
	}

	if err := db.UpdateUser(user); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "Password updated successfully")
	}
}
