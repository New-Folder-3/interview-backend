package db

import (
	"github.com/pkg/errors"
	"interview-backend/internal/model"
	"time"
)

func CreateToken(t *model.Token) error {
	return errors.WithStack(db.Create(t).Error)
}

func GetToken(token string) (*model.Token, error) {
	var t model.Token
	ClearExpiredTokens()
	err := db.Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &t, nil
}

func DeleteToken(token string) error {
	return db.Where("token = ?", token).Delete(&model.Token{}).Error
}

func ClearExpiredTokens() error {
	return errors.WithStack(db.Where("expire_ts < ?", time.Now().Unix()).Delete(&model.Token{}).Error)
}
