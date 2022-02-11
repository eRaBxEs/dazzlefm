package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetContact ...
func (s Office) GetContact(id int64) (form form.Contact, err error) {
	record := model.Contact{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Take(&record).Error
	})

	err = copier.Copy(&form, record)

	return
}

// ListContact ...
func (s Office) ListContact(pageNo, offset int) (records []form.Contact, err error) {
	records = []form.Contact{}
	// err = s.dbc.Order("date desc").Find(&records).Error

	if pageNo == 0 {

	}
	if pageNo == 1 {
		err = s.dbc.Order("date desc").Limit(offset).Find(&records).Error
	}
	if pageNo > 1 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("date desc").Limit(offset).Find(&records).Error
	}

	return
}

// CountContact ...
func (s Office) CountContact() (cnt int, err error) {

	record := form.Contact{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return

}

// AddContact ...
func (s Office) AddContact(form *form.Contact) (err error) {
	record := &model.Contact{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.Date = time.Now()

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SaveContact ...
func (s Office) SaveContact(id int64, form *form.Contact) (err error) {
	record := &model.Contact{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteContact ...
func (s Office) DeleteContact(id int64) (err error) {
	record := model.Contact{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
