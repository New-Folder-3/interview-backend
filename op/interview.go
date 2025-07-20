package op

import (
	"fmt"
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
	"time"
)

func NewInterview(Username string) (string, error) {
	interview := model.NewDefaultInterview()
	interview.Username = Username
	interview.CreateAt = time.Now().Unix()
	userDB, err := db.GetUser(Username)
	if err != nil {
		return "", errors.WithStack(err)
	}
	userDB.Interviews = util.StringListToDB(append(util.DBToStringList(userDB.Interviews), interview.ID))
	if err = db.UpdateUser(userDB); err != nil {
		return "", errors.WithStack(err)
	}
	if err = db.CreateInterview(&interview); err != nil {
		return "", errors.WithStack(err)
	}
	return interview.ID, nil
}

func DeleteInterview(InterviewID string) error {
	interview, err := db.GetInterview(InterviewID)
	if err != nil {
		return errors.WithStack(err)
	}
	userDB, err := db.GetUser(interview.Username)
	if err != nil {
		return errors.WithStack(err)
	}
	userDB.Interviews = util.RemoveFromDBList(userDB.Interviews, interview.ID)
	if err = db.ReplaceUser(userDB); err != nil {
		return errors.WithStack(err)
	}
	if err = db.DeleteInterview(InterviewID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetModelComment(InterviewID string) (*model.Comment, error) {
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret, err := db.GetComment(interviewDB.ModelComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ret, nil
}

func GetUserComment(InterviewID string) (*[]model.Comment, error) {
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	comments := util.DBToStringList(interviewDB.UserComments)
	var ret []model.Comment
	for _, comment := range comments {
		c, err := db.GetComment(comment)
		if err != nil {
			util.ErrorPrinter(errors.WithStack(err))
			continue
		}
		ret = append(ret, *c)
	}
	return &ret, nil
}

func AddModelComment(InterviewID string, Comment model.CommentContent) error {
	commentDB := model.Comment{
		ID:             util.GenerateToken(16),
		CommentContent: Comment,
	}
	if err := db.CreateComment(&commentDB); err != nil {
		return errors.WithStack(err)
	}
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return errors.WithStack(err)
	}
	interviewDB.ModelComments = commentDB.ID
	if err = db.UpdateInterview(interviewDB); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AddUserComment(InterviewID string, Comment model.CommentContent) error {
	commentDB := model.Comment{
		ID:             util.GenerateToken(16),
		CommentContent: Comment,
	}
	if err := db.CreateComment(&commentDB); err != nil {
		return errors.WithStack(err)
	}
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return errors.WithStack(err)
	}
	interviewDB.UserComments = util.StringListToDB(append(util.DBToStringList(interviewDB.UserComments), commentDB.ID))
	if err = db.UpdateInterview(interviewDB); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetDimension(ConversationID, InterviewID, Username string) (model.InterviewDimension, error) {
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return model.InterviewDimension{}, errors.WithStack(err)
	}

	if ConversationID == "" {
		ConversationID = interviewDB.ResultConversation
	}

	userDB, err := db.GetUser(Username)
	if err != nil {
		return model.InterviewDimension{}, errors.WithStack(err)
	}

	response, err := FastChatTxt(ConversationID, fmt.Sprintf(conf.PromptTemplate[conf.DimensionPrompt], userDB.Job))
	if err != nil {
		return model.InterviewDimension{}, errors.WithStack(err)
	}

	dimensions := util.ChatStringToFloatSlice(response)
	for index, dimension := range dimensions {
		switch index {
		case 0:
			interviewDB.InterviewDimension.Hard = dimension
		case 1:
			interviewDB.InterviewDimension.Soft = dimension
		case 2:
			interviewDB.InterviewDimension.Potential = dimension
		case 3:
			interviewDB.InterviewDimension.Confidence = dimension
		case 4:
			interviewDB.InterviewDimension.Development = dimension
		case 5:
			interviewDB.InterviewDimension.Fit = dimension
		}
	}

	if err = db.UpdateInterview(interviewDB); err != nil {
		return model.InterviewDimension{}, errors.WithStack(err)
	}
	return interviewDB.InterviewDimension, nil
}

func GetKeywords(ConversationID, InterviewID string) ([]string, error) {
	interviewDB, err := db.GetInterview(InterviewID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if ConversationID == "" {
		ConversationID = interviewDB.ResultConversation
	}

	response, err := FastChatTxt(ConversationID,
		fmt.Sprintf(conf.PromptTemplate[conf.KeywordsPrompt], conf.Keywords))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	keywords := util.ChatStringToStringSlice(response)
	interviewDB.Keywords = util.StringListToDB(keywords)

	if err = db.UpdateInterview(interviewDB); err != nil {
		return nil, errors.WithStack(err)
	}
	return keywords, nil
}
