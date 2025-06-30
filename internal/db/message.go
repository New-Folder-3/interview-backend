package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func GetMessage(MessageID string) (*model.Message, error) {
	var talking model.Message
	if err := db.Where("id = ?", MessageID).First(&talking).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &talking, nil
}

func CreateMessage(t *model.Message) error {
	return errors.WithStack(db.Create(t).Error)
}

func DeleteMessage(MessageID string) error {
	return errors.WithStack(db.Where("id = ?", MessageID).Delete(&model.Message{}).Error)
}

func UpdateMessage(t *model.Message) error {
	return errors.WithStack(db.Model(t).Updates(t).Error)
}

func ReplaceMessage(t *model.Message) error {
	return errors.WithStack(db.Model(t).Save(t).Error)
}
