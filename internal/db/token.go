package db

import (
	"interview-backend/internal/model"
	"time"
)

func CreateToken(t *model.Token) error {
	return db.Create(t).Error
}

func GetToken(token string) (*model.Token, error) {
	var t model.Token
	ClearExpiredTokens()
	err := db.Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func DeleteToken(user, token string) error {
	return db.Where("user_id = ? AND token = ?", user, token).Delete(&model.Token{}).Error
}

func ClearExpiredTokens() error {
	return db.Where("expire_ts < ?", time.Now().Unix()).Delete(&model.Token{}).Error
}
