package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetRole ...
func (s Office) GetRole(id int64) (form form.PersonRole, err error) {
	record := model.PersonRole{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// ListRole ...
func (s Office) ListRole() (records []form.PersonRole, err error) {
	records = []form.PersonRole{}
	err = s.dbc.Order("id").Find(&records).Error

	return
}

// AddRole ...
func (s Office) AddRole(form *form.PersonRole) (err error) {
	record := &model.PersonRole{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})
	if err != nil {
		return
	}

	form.ID = record.ID

	return
}

// SaveRole ...
func (s Office) SaveRole(id int64, form *form.PersonRole) (err error) {
	record := &model.PersonRole{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteRole ...
func (s Office) DeleteRole(id int64) (err error) {
	record := model.PersonRole{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		count := struct {
			Count int
		}{}
		err := tx.Raw(`
			select 
				count(id)
			from
				person
			where
				role = ?
		`, id).Scan(&count).Error
		if err != nil {
			s.log.Error(err)
			return err
		}

		if count.Count > 0 {
			return fmt.Errorf("role is used by %d person", count.Count)
		}

		if err = tx.Delete(record).Error; err != nil {
			s.log.Error(err)
			return err
		}

		return nil
	})

	return
}
