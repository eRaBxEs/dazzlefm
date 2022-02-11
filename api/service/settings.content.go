package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// ListSettings ...
func (s Content) ListSettings() (record form.Settings, err error) {
	record = form.Settings{}
	err = s.dbc.Take(&record).Error

	return
}

// AddSettings ...
func (s Content) AddSettings(form *form.Settings) (err error) {

	record := &model.Settings{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	// record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Model(record).Update(record).Error
	})

	//form.ID = record.ID

	return
}
