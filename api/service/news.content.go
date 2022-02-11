package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetNews ...
func (s Content) GetNews(id int64) (form form.News, err error) {
	record := model.News{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// CountNews ...
func (s Content) CountNews() (cnt int, err error) {
	record := form.News{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return
}

// ListNews ...
func (s Content) ListNews(pageNo, offset int) (records []form.News, err error) {
	records = []form.News{}

	if pageNo == 0 {

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

// AddNews ...
func (s Content) AddNews(form *form.News) (err error) {
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
	// s.log.Debugf("%d", image.Size)

	record := &model.News{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SaveNews ...
func (s Content) SaveNews(id int64, form *form.News) (err error) {
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

	record := &model.News{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteNews ...
func (s Content) DeleteNews(id int64) (err error) {
	record := model.News{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
