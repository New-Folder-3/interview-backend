package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	modelComment, err1 := op.GetModelComment(request.InterviewID)
	userComment, err2 := op.GetUserComment(request.InterviewID)

	var msg string
	if err1 != nil {
		msg += "Get Model Comment Failed: " + err1.Error() + "; "
	}
	if err2 != nil {
		msg += "Get User Comment Failed: " + err2.Error() + "; "
	}
	util.SuccessResp(c, map[string]interface{}{
		"model_comment": modelComment,
		"user_comment":  userComment,
	}, msg)
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
		err = op.AddModelComment(request.InterviewID, request.CommentContent)
	case "user":
		err = op.AddUserComment(request.InterviewID, request.CommentContent)
	default:
		err = errors.New("invalid comment role")
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
