package db

import (
	"github.com/pkg/errors"
	"interview-backend/internal/model"
)

func CreateContent(c *model.Content) error {
	return errors.WithStack(db.Create(c).Error)
}

func UpdateContent(c *model.Content) error {
	return errors.WithStack(db.Model(c).Updates(c).Error)
}

func DeleteContent(c *model.Content) error {
	return errors.WithStack(db.Delete(c).Error)
}

func GetContent(ContentID string) (*model.Content, error) {
	var content model.Content
	if err := db.Where("id = ?", ContentID).First(&content).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &content, nil
}
