package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetSnippet ...
func (s Content) GetSnippet(id int64) (form form.Script, err error) {
	record := model.Script{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// ListSnippet ...
func (s Content) ListSnippet() (records []form.Script, err error) {
	records = []form.Script{}
	err = s.dbc.Order("id desc").Find(&records).Error

	return
}

// AddSnippet ...
func (s Content) AddSnippet(form *form.Script) (err error) {

	record := &model.Script{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SaveSnippet ...
func (s Content) SaveSnippet(id int64, form *form.Script) (err error) {

	record := &model.Script{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteSnippet ...
func (s Content) DeleteSnippet(id int64) (err error) {
	record := model.Script{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
