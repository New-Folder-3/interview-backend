package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func CreateContent(c *model.Content) error {
	return errors.WithStack(db.Create(c).Error)
}

func UpdateContent(c *model.Content) error {
	return errors.WithStack(db.Model(c).Updates(c).Error)
}

func ReplaceContent(c *model.Content) error {
	return errors.WithStack(db.Model(c).Save(c).Error)
}

func DeleteContent(ID string) error {
	return errors.WithStack(db.Where("id = ?", ID).Delete(&model.Content{}).Error)
}

func GetContent(ContentID string) (*model.Content, error) {
	var content model.Content
	if err := db.Where("id = ?", ContentID).First(&content).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &content, nil
}
