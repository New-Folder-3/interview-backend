package handles

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/db"
	"interview-backend/util"
	"time"
)

func GetUser(c *gin.Context) {
	var getRequest UserRequest
	if err := c.ShouldBind(&getRequest); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	user, err := db.GetUser(getRequest.Username)
	if err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "User retrieved successfully")
	}
}

func UpdateUser(c *gin.Context) {
	var updateRequest UserRequest
	if err := c.ShouldBind(&updateRequest); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}
	user, err := db.GetUser(updateRequest.Username)
	if err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
		return
	}

	user.Email = updateRequest.Email
	user.Phone = updateRequest.Phone
	user.Name = updateRequest.Nickname

	if updateRequest.OldPwd == user.PwdHash {
		user.PwdHash = updateRequest.NewPwd
		user.PwdTS = time.Now().Unix()
	} else if updateRequest.OldPwd != "" {
		util.ErrorResp(c, "Incorrect Old Password, Password changed failed", 400, false)
	}

	if err := db.UpdateUser(user); err != nil {
		util.ErrorResp(c, err.Error(), 400, false)
	} else {
		util.SuccessResp(c, user, "User updated successfully")
	}
}
