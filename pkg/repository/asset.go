package repository

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
)

func FindAllAsset(o *[]model.Asset) (err error) {
	if err = config.DB.Find(o).Error; err != nil {
		return err
	}
	return nil
}

func FindPageAsset(pageIndex, pageSize int, total *int64, o *[]model.Asset, name string) (err error) {
	db := config.DB
	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if err := db.Find(&o).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Count(total).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewAsset(o *model.Asset) (err error) {
	if err = config.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindAssetById(o *model.Asset, id string) (err error) {

	if err := config.DB.Where("id = ?", id).First(o).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAssetById(o *model.Asset, id string) (err error) {
	o.ID = id
	config.DB.Updates(o)
	return nil
}

func DeleteAssetById(id string) (err error) {
	config.DB.Delete(&model.Asset{}, id)
	return nil
}
