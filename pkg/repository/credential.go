package repository

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
)

func FindAllCredential(o *[]model.Credential) (err error) {
	if err = config.DB.Find(o).Error; err != nil {
		return err
	}
	return nil
}

func FindPageCredential(pageIndex, pageSize int, total *int64, o *[]model.Credential, name string) (err error) {
	db := config.DB
	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if err := db.Order("created desc").Find(&o).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Count(total).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewCredential(o *model.Credential) (err error) {
	if err = config.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindCredentialById(o *model.Credential, id string) (err error) {

	if err := config.DB.Where("id = ?", id).First(o).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCredentialById(o *model.Credential, id string) (err error) {
	o.ID = id
	config.DB.Updates(o)
	return nil
}

func DeleteCredentialById(id string) (err error) {
	config.DB.Delete(&model.Credential{}, id)
	return nil
}
