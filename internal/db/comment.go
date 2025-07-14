package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func GetComment(CommentID string) (*model.Comment, error) {
	var comment model.Comment
	if err := db.Where("id = ?", CommentID).First(&comment).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &comment, nil
}

func CreateComment(i *model.Comment) error {
	return errors.WithStack(db.Create(i).Error)
}

func UpdateComment(i *model.Comment) error {
	return errors.WithStack(db.Model(i).Updates(i).Error)
}

func DeleteComment(CommentID string) error {
	return errors.WithStack(db.Where("id = ?", CommentID).Delete(&model.Comment{}).Error)
}

func ReplaceComment(i *model.Comment) error {
	return errors.WithStack(db.Model(i).Save(i).Error)
}
