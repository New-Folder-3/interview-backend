package db

import (
	"github.com/pkg/errors"
	"interview/internal/model"
)

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(&u).Error)
}

func UpdateUser(u *model.User) error {
	return errors.WithStack(db.Model(u).Updates(u).Error)
}

func ReplaceUser(u *model.User) error {
	return errors.WithStack(db.Model(u).Save(u).Error)
}

func DeleteUser(userID string) error {
	return errors.WithStack(db.Where("id = ?", userID).Delete(&model.User{}).Error)
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
