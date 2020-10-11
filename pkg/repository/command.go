package repository

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
)

func FindAllCommand(o *[]model.Command) (err error) {
	if err = config.DB.Find(o).Error; err != nil {
		return err
	}
	return nil
}

func FindPageCommand(pageIndex, pageSize int, total *int64, o *[]model.Command, name, content string) (err error) {

	db := config.DB
	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if len(content) > 0 {
		db = db.Where("content like ?", "%"+content+"%")
	}

	if err := db.Find(&o).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Count(total).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewCommand(o *model.Command) (err error) {
	if err = config.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindCommandById(o *model.Command, id string) (err error) {

	if err := config.DB.Where("id = ?", id).First(o).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCommandById(o *model.Command, id string) (err error) {
	o.ID = id
	config.DB.Updates(o)
	return nil
}

func DeleteCommandById(id string) (err error) {
	config.DB.Delete(&model.Command{}, id)
	return nil
}
