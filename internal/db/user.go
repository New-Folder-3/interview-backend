package db

import (
	"github.com/pkg/errors"
	"interview-backend/internal/model"
)

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(&u).Error)
}

func UpdateUser(u *model.User) error {
	return db.Model(u).Updates(u).Error
}

func DeleteUser(u *model.User) error {
	return errors.WithStack(db.Delete(u).Error)
}

func GetUser(info string) (*model.User, error) {
	var user model.User
	if err := db.Where("email = ? OR phone = ? OR id = ?", info, info, info).First(&user).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

func DeleteUserByUUID(uuid string) error {
	return errors.WithStack(db.Delete(model.User{}, uuid).Error)
}
