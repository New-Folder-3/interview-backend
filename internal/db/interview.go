package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func GetInterview(InterviewID string) (*model.Interview, error) {
	var interview model.Interview
	if err := db.Where("id = ?", InterviewID).First(&interview).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &interview, nil
}

func CreateInterview(i *model.Interview) error {
	return errors.WithStack(db.Create(i).Error)
}

func DeleteInterview(InterviewID string) error {
	return errors.WithStack(db.Where("id = ?", InterviewID).Delete(&model.Interview{}).Error)
}

func UpdateInterview(i *model.Interview) error {
	return errors.WithStack(db.Model(i).Updates(i).Error)
}

func ReplaceInterview(i *model.Interview) error {
	return errors.WithStack(db.Model(i).Save(i).Error)
}
