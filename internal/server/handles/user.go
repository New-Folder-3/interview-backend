package handles

import (
	"github.com/gin-gonic/gin"
	"interview/internal/db"
	"interview/util"
	"time"
)

func GetUser(c *gin.Context) {
	var getRequest UserRequest
	if err := c.ShouldBind(&getRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
		return
	}
	userShould, _ := c.Get("user")
	if getRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	user, err := db.GetUser(getRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
	} else {
		util.SuccessResp(c, user, "User retrieved successfully")
	}
}

func UpdateUserInfo(c *gin.Context) {
	var updateRequest UserRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
		return
	}
	userShould, _ := c.Get("user")
	if updateRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	user, err := db.GetUser(updateRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
		return
	}
	user.UserChangable = updateRequest.UserChangable

	if err := db.UpdateUser(user); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
	} else {
		util.SuccessResp(c, user, "User updated successfully")
	}
}

func UpdateUserPwd(c *gin.Context) {
	var updatePwdRequest UserRequest
	if err := c.ShouldBindJSON(&updatePwdRequest); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
		return
	}
	userShould, _ := c.Get("user")
	if updatePwdRequest.Username != userShould {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	user, err := db.GetUser(updatePwdRequest.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
		return
	}

	if updatePwdRequest.PwdHash == user.PwdHash {
		user.PwdHash = updatePwdRequest.NewPwdHash
		user.PwdTS = time.Now().Unix()
	} else if updatePwdRequest.PwdHash != "" {
		util.ErrorResp(c, "Incorrect Old Password, Password changed failed", 400)
		return
	}

	if err := db.UpdateUser(user); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 400)
	} else {
		util.SuccessResp(c, user, "Password updated successfully")
	}
}
