package db

import (
	"github.com/pkg/errors"
	"interview-backend/internal/model"
	"interview-backend/util"
)

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(&u).Error)
}

func UpdateUser(u *model.User) error {
	return errors.WithStack(db.Model(&u).Error)
}

func DeleteUser(u *model.User) error {
	return errors.WithStack(db.Delete(u).Error)
}

func GetUsersByPage(pageIndex, pageSize int) ([]model.User, int64, error) {
	var count int64 // 总数
	var ret []model.User
	userDB := db.Model(&model.User{})
	if err := userDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed to count users")
	}
	if err := userDB.Order(util.DbGetColumnName("uuid")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&ret).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed to find users")
	}
	return ret, count, nil
}

func DeleteUserByUUID(uuid string) error {
	return errors.WithStack(db.Delete(model.User{}, uuid).Error)
}
