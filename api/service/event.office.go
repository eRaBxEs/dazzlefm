package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetEvent ...
func (s Office) GetEvent(id int64) (form form.Event, err error) {
	record := model.Event{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// CountEvent ...
func (s Office) CountEvent() (cnt int, err error) {

	record := form.Event{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return

}

// ListEvent ...
func (s Office) ListEvent(pageNo, offset int) (records []form.Event, err error) {
	records = []form.Event{}

	// err = s.dbc.Order("date asc").Find(&records).Error

	if pageNo == 0 {

	}
	if pageNo == 1 {
		err = s.dbc.Order("date asc").Limit(offset).Find(&records).Error
	}
	if pageNo > 1 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("date asc").Limit(offset).Find(&records).Error
	}

	return
}

// AddEvent ...
func (s Office) AddEvent(frm *form.Event) (err error) {

	record := &model.Event{}
	if err = copier.Copy(record, frm); err != nil {
		s.log.Error(err)
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})
	if err != nil {
		s.log.Error(err)
		return
	}

	frm.ID = record.ID

	return
}

// SaveEvent ...
func (s Office) SaveEvent(id int64, form *form.Event) (err error) {

	record := &model.Event{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteEvent ...
func (s Office) DeleteEvent(id int64) (err error) {
	record := model.Event{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}

// GetLiveEvent ...
func (s Office) GetLiveEvent() (form form.Event, err error) {
	record := model.Event{}
	todayDate, err := util.GetTodayDate()
	if err != nil {
		return
	}

	err = s.dbc.Raw("SELECT * FROM event WHERE (date >= ?) ORDER BY date ASC LIMIT 1 ", todayDate).Scan(&record).Error
	if err != nil {
		return
	}

	err = copier.Copy(&form, record)
	if err != nil {
		return
	}

	return
}

// GetOtherEvents ...
func (s Office) GetOtherEvents() (records []form.Event, err error) {
	records = []form.Event{}
	todayDate, err := util.GetTodayDate()
	if err != nil {
		return records, err
	}

	err = s.dbc.Raw("SELECT * FROM event  WHERE (date >= ?) ORDER BY date ASC OFFSET 1 LIMIT 5 ", todayDate).Scan(&records).Error
	if err != nil {
		return records, err
	}

	return
}
