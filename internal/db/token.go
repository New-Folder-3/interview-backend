package db

import (
	"interview-backend/internal/model"
	"time"
)

func CreateToken(t *model.Token) error {
	t.ExpireTS = time.Now().Unix()
	return db.Create(t).Error
}

func GetToken(user, token string) (*model.Token, error) {
	var t model.Token
	err := db.Where("user_id = ? AND token = ?", user, token).First(&t).Error
	if err != nil {
		return nil, err
	}
	if t.ExpireTS < time.Now().Unix() {
		DeleteToken(user, token)
		return nil, nil // Token expired
	}
	return &t, nil
}

func DeleteToken(user, token string) error {
	return db.Where("user_id = ? AND token = ?", user, token).Delete(&model.Token{}).Error
}
