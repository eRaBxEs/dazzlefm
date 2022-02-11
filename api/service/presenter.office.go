package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetPresenter ...
func (s Office) GetPresenter(id int64) (form form.Presenter, err error) {
	record := model.Presenter{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Take(&record).Error
	})

	err = copier.Copy(&form, record)

	return
}

// CountPresenter ...
func (s Office) CountPresenter() (cnt int, err error) {

	record := form.Presenter{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return

}

// ListPresenter ...
func (s Office) ListPresenter(pageNo, offset, status int) (records []form.Presenter, err error) {
	records = []form.Presenter{}
	if pageNo == 0 && status == 0 {
		err = s.dbc.Order("name").Find(&records).Error
	}
	if pageNo == 0 && status != 0 {
		err = s.dbc.Order("name").Where("status = ?", status).Find(&records).Error
	}
	if pageNo == 1 && status == 0 {
		err = s.dbc.Order("name").Limit(offset).Find(&records).Error
	}
	if pageNo == 1 && status != 0 {
		err = s.dbc.Order("name").Limit(offset).Where("status = ?", status).Find(&records).Error
	}
	if pageNo > 1 && status == 0 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("name").Limit(offset).Find(&records).Error
	}
	if pageNo > 1 && status != 0 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("name").Limit(offset).Where("status = ?", status).Find(&records).Error
	}

	return
}

// AddPresenter ...
func (s Office) AddPresenter(form *form.Presenter) (err error) {
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

	record := &model.Presenter{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SavePresenter ...
func (s Office) SavePresenter(id int64, form *form.Presenter) (err error) {
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

	record := &model.Presenter{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeletePresenter ...
func (s Office) DeletePresenter(id int64) (err error) {
	record := model.Presenter{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
