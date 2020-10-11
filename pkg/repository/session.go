package repository

import (
	"dogo/pkg/config"
	"dogo/pkg/model"
)

func FindAllSession(o *[]model.Session) (err error) {
	if err = config.DB.Find(o).Error; err != nil {
		return err
	}
	return nil
}

func FindPageSession(pageIndex, pageSize int, total *int64, o *[]model.Session, status string) (err error) {

	db := config.DB
	if len(status) > 0 {
		db = db.Where("status = ?", status)
	}

	if err := db.Find(&o).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Count(total).Error; err != nil {
		return err
	}
	return nil
}

func CreateNewSession(o *model.Session) (err error) {
	if err = config.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindSessionById(o *model.Session, id string) (err error) {
	if err := config.DB.Where("id = ?", id).First(o).Error; err != nil {
		return err
	}
	return nil
}

func FindSessionByConnectionId(o *model.Session, connectionId string) (err error) {
	if err := config.DB.Where("connection_id = ?", connectionId).First(o).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSessionById(o *model.Session, id string) (err error) {
	o.ID = id
	config.DB.Updates(o)
	return nil
}

func DeleteSessionById(id string) (err error) {
	config.DB.Delete(&model.Session{}, id)
	return nil
}
