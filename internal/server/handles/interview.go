package handles

import (
	"github.com/gin-gonic/gin"
	"interview/internal/db"
	"interview/op"
	"interview/util"
)

func NewInterview(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	interviewID, err := op.NewInterview(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, interviewID, "Create Interview Success")
}

func GetInterview(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}
	interview, err := db.GetInterview(request.InterviewID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, interview, "Get Interview Success")
}

func UpdateInterview(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	interview, err := db.GetInterview(request.InterviewID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	interview.InterviewChangeable = request.InterviewChangeable
	if err = db.UpdateInterview(interview); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, interview, "Update Interview Success")
}

func DeleteInterview(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	if err := op.DeleteInterview(request.InterviewID); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	util.SuccessResp(c, nil, "Delete Interview Success")
}

func GetAllInterviews(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	userDB, err := db.GetUser(request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	interviews := util.DBToStringList(userDB.Interviews)
	util.SuccessResp(c, interviews, "Get All Interviews Successfully")
}

func GetComment(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	modelComment, err := op.GetModelComment(request.InterviewID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	userComment, err := op.GetUserComment(request.InterviewID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}
	util.SuccessResp(c, map[string]interface{}{
		"model_comment": modelComment,
		"user_comment":  userComment,
	}, "Get Comment Success")
}

func AddComment(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	var err error
	switch request.CommentRole {
	case "model":
		err = op.AddModelComment(request.InterviewID, request.Comment)
	case "user":
		err = op.AddUserComment(request.InterviewID, request.Comment)
	}
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	util.SuccessResp(c, nil, "Add Comment Successfully")
}

func GetDimension(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	interviewDimension, err := op.GetDimension(request.ConversationID, request.InterviewID, request.Username)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	util.SuccessResp(c, interviewDimension, "Get Dimension Successfully")
}

func GetKeywords(c *gin.Context) {
	var request InterviewRequest
	if err := c.ShouldBind(&request); err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, "Invalid Request", 400)
		return
	}
	userToken, _ := c.Get("user")
	if userToken.(string) != request.Username {
		util.ErrorResp(c, "Username mismatch", 400)
		return
	}

	keywords, err := op.GetKeywords(request.ConversationID, request.InterviewID)
	if err != nil {
		util.ErrorPrinter(err)
		util.ErrorResp(c, err.Error(), 500)
		return
	}

	util.SuccessResp(c, keywords, "Get Keywords Successfully")
}
