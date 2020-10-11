package repository

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
)

func FindAllUser(o *[]model.User) (err error) {
	if err = config.DB.Find(o).Error; err != nil {
		return err
	}
	return nil
}

func FindPageUser(pageIndex, pageSize int, total *int64, o *[]model.User) (err error) {
	if err := config.DB.Find(&o).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Count(total).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewUser(o *model.User) (err error) {
	if err = config.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindUserById(o *model.User, id string) (err error) {
	if err := config.DB.Where("id = ?", id).First(o).Error; err != nil {
		return err
	}
	return nil
}

func FindUserByUsername(o *model.User, username string) (err error) {
	if err := config.DB.Where("username = ?", username).First(o).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUserById(o *model.User, id string) (err error) {
	o.ID = id
	config.DB.Updates(o)
	return nil
}

func DeleteUserById(id string) (err error) {
	config.DB.Where("id = ?", id).Delete(&model.User{})
	return nil
}
