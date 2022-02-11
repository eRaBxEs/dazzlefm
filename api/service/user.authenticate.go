package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// FindUserByUserName ...
func (s Authenticate) FindUserByUserName(user string) (form form.User, err error) {
	record := model.User{}
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Take(&record, "user_name = ?", user).Error
	})
	if err != nil {
		return
	}

	if err = copier.Copy(&form, record); err != nil {
		return
	}

	return
}

// GetUser ...
func (s Authenticate) GetUser(id int64) (form form.User, err error) {
	record := model.User{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Take(&record).Error
	})
	if err != nil {
		return
	}

	if err = copier.Copy(&form, record); err != nil {
		return
	}

	form.Password = "***"

	return
}

// ListUser ...
func (s Authenticate) ListUser() (records []form.User, err error) {
	records = []form.User{}
	if err = s.dbc.Order("first_name").Find(&records).Error; err != nil {
		return
	}

	for i := 0; i < len(records); i++ {
		records[i].Password = "***"
	}

	return
}

// AddUser ...
func (s Authenticate) AddUser(form *form.User) (err error) {
	record := &model.User{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	if record.Password != "***" || len(record.Password) > 0 {
		retv, err := util.HashPassword(record.Password)
		if err != nil {
			return err
		}

		record.Password = retv
		form.Password = "***"
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

// SaveUser ...
func (s Authenticate) SaveUser(id int64, form *form.User) (err error) {
	record := &model.User{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	if record.Password != "***" || len(record.Password) > 0 {
		retv, err := util.HashPassword(record.Password)
		if err != nil {
			return err
		}

		record.Password = retv
		form.Password = "***"
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteUser ...
func (s Authenticate) DeleteUser(id int64) (err error) {
	record := model.User{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}
