package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetGallery ...
func (s Content) GetGallery(id int64) (form form.Gallery, err error) {
	record := model.Gallery{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// ListGallery ...
func (s Content) ListGallery(pageNo, offset int) (records []form.Gallery, err error) {
	records = []form.Gallery{}

	if pageNo == 0 {
		err = s.dbc.Order("id desc").Find(&records).Error
	}

	if pageNo == 1 {
		err = s.dbc.Order("id desc").Limit(offset).Find(&records).Error
	}
	if pageNo > 1 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("id desc").Limit(offset).Find(&records).Error
	}

	return
}

// CountGallery ...
func (s Content) CountGallery() (cnt int, err error) {

	record := form.Gallery{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return
}

// AddGallery ...
func (s Content) AddGallery(form *form.Gallery) (err error) {
	image, err := util.SaveImageData(s.env, form.Image)
	if err != nil {
		s.log.Error(err)
		return err
	}

	form.Image, err = json.Marshal(image)
	if err != nil {
		s.log.Error(err)
		return err
	}

	record := &model.Gallery{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SaveGallery ...
func (s Content) SaveGallery(id int64, form *form.Gallery) (err error) {
	image, err := util.SaveImageData(s.env, form.Image)
	if err != nil {
		s.log.Error(err)
		return err
	}

	form.Image, err = json.Marshal(image)
	if err != nil {
		s.log.Error(err)
		return err
	}

	record := &model.Gallery{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteGallery ...
func (s Content) DeleteGallery(id int64) (err error) {
	record := model.Gallery{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
