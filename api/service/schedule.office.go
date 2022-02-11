package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetSchedule ...
func (s Office) GetSchedule(id int64) (form form.Schedule, err error) {
	record := model.Schedule{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	record.StartTime = record.StartTime.AddDate(1970, 1, 1)
	record.EndTime = record.EndTime.AddDate(1970, 1, 1)
	err = copier.Copy(&form, record)

	return
}

// CountSchedule ...
func (s Office) CountSchedule() (cnt int, err error) {

	record := form.Schedule{}
	err = s.dbc.Model(&record).Count(&cnt).Error

	return

}

// ListSchedule ...
func (s Office) ListSchedule(pageNo, offset, day int) (records []form.Schedule, err error) {
	records = []form.Schedule{}

	// err = s.dbc.Order("id").Find(&records).Error

	if pageNo == 0 && offset == 0 && day == 0 {

	}

	if pageNo == 1 && day == 0 {
		err = s.dbc.Order("id desc").Limit(offset).Find(&records).Error
	}
	if pageNo > 1 && day == 0 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("id desc").Limit(offset).Find(&records).Error
	}

	if pageNo == 1 && offset > 0 && day > 0 {

		err = s.dbc.Order("start_time asc").Limit(offset).Find(&records, "day = ?", day).Error
	}

	if pageNo > 1 && offset > 0 && day > 0 {
		page := (pageNo - 1) * offset
		err = s.dbc.Offset(page).Order("start_time asc").Limit(offset).Find(&records, "day = ?", day).Error
	}

	for i := 0; i < len(records); i++ {
		records[i].StartTime = records[i].StartTime.AddDate(1970, 1, 1)
		records[i].EndTime = records[i].EndTime.AddDate(1970, 1, 1)
	}

	return
}

// ListASchedule ...
func (s Office) ListASchedule(presenterID int64) (records []form.Schedule, err error) {
	records = []form.Schedule{}
	err = s.dbc.Order("day asc").Find(&records, "presenter_id = ? OR copresenter_id = ?", presenterID, presenterID).Error

	for i := 0; i < len(records); i++ {
		records[i].StartTime = records[i].StartTime.AddDate(1970, 1, 1)
		records[i].EndTime = records[i].EndTime.AddDate(1970, 1, 1)
	}

	return
}

// AddSchedule ...
func (s Office) AddSchedule(frm *form.Schedule) (ifExist bool, form form.Schedule, err error) {

	record := &model.Schedule{}
	if err = copier.Copy(record, frm); err != nil {
		s.log.Error(err)
		return
	}

	err = s.dbc.Raw(`SELECT * FROM schedule WHERE day = ? AND  ? >= start_time 
	AND ? <= end_time ORDER BY start_time LIMIT 1`, record.Day, record.StartTime, record.StartTime).Scan(&record).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// To check if the record exists
	if record.ID > 0 {
		// Set ifExist to true and return
		ifExist = true
		err = copier.Copy(&form, record)
		if err != nil {
			return
		}
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

// SaveSchedule ...
func (s Office) SaveSchedule(id int64, form *form.Schedule) (ifExist bool, frm form.Schedule, err error) {

	record := &model.Schedule{}

	err = s.dbc.Raw(`SELECT * FROM schedule WHERE day = ? AND  ? >= start_time 
	AND ? <= end_time ORDER BY start_time LIMIT 1`, form.Day, form.StartTime, form.StartTime).Scan(&record).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// To check if the record exists
	if record.ID > 0 && record.ID != id {
		// Set ifExist to true and return
		ifExist = true
		err = copier.Copy(&frm, record)
		if err != nil {
			return
		}
		return
	}

	// over writing the content of record got from the database

	err = copier.Copy(&record, form)
	if err != nil {
		return
	}
	record.ID = id // resetting the id of record
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	if err != nil {
		s.log.Error(err)
		return
	}

	return
}

// DeleteSchedule ...
func (s Office) DeleteSchedule(id int64) (err error) {
	record := model.Schedule{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}

// GetLiveSchedule ...
func (s Office) GetLiveSchedule(dayVal int, timeVal string) (form form.Schedule, err error) {

	record := model.Schedule{}

	err = s.dbc.Raw(`
		SELECT * FROM schedule WHERE (day = ? start_time <= ? AND end_time > ?) ORDER BY start_time
	`, dayVal, timeVal, timeVal).Scan(&record).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	record.StartTime = record.StartTime.AddDate(1970, 1, 1)
	record.EndTime = record.EndTime.AddDate(1970, 1, 1)

	err = copier.Copy(&form, record)
	if err != nil {
		return
	}

	return
}

// GetScheduleSummary ...
func (s Office) GetScheduleSummary(dayVal int, timeVal string) (records []form.ScheduleHighlight, err error) {

	records = []form.ScheduleHighlight{}

	qry := s.dbc.Table("schedule as s").Select(`
		s.*,
		p.*,
		c.name as co_name,
		c.image as co_image`).
		Where(
			"s.day = ?", dayVal)

	qry = qry.Where("(? between s.start_time AND s.end_time) OR (s.day = ? AND end_time > ?)", timeVal, dayVal, timeVal).
		// qry = qry.Where("(? between s.start_time AND s.end_time) union select * from schedule where s.day= ? AND end_time > ?",timeVal, dayVal, timeVal).
		Joins("left join presenter as p on p.id = s.presenter_id").
		Joins("left join presenter as c on c.id = s.copresenter_id").
		Order("start_time").Limit(5)
	err = qry.Scan(&records).Error

	// err = s.dbc.Raw(`SELECT
	// 	s.*,
	// 	p.*,
	// 	c.name as co_name,
	// 	c.image as co_image FROM schedule as s
	// 		left join presenter as p on p.id = s.presenter_id
	// 		left join presenter as c on c.id = s.copresenter_id
	// 		WHERE (s.day = '2') AND (('22:49:28' between s.start_time AND s.end_time) union select * from schedule where s.day= '2' AND end_time > '22:49:28') ORDER BY start_time LIMIT 5

	// `)

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if len(records) > 0 {
		// add date to make StartTime and EndTime legitimate timestamp structures
		for i := 0; i < len(records); i++ {
			records[i].StartTime = records[i].StartTime.AddDate(1970, 1, 1)
			records[i].EndTime = records[i].EndTime.AddDate(1970, 1, 1)
		}
	}

	return

}

// GetAllSchedules ...
func (s Office) GetAllSchedules(dayVal int) (records []form.ScheduleHighlight, err error) {

	records = []form.ScheduleHighlight{}

	qry := s.dbc.Table("schedule as s").Select(`
		s.*,
		p.name,
		c.name as co_name`).Where("day = ?", dayVal).
		Joins("left join presenter as p on p.id = s.presenter_id").
		Joins("left join presenter as c on c.id = s.copresenter_id").
		Order("start_time")
	err = qry.Scan(&records).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	for i := 0; i < len(records); i++ {
		records[i].StartTime = records[i].StartTime.AddDate(1970, 1, 1)
		records[i].EndTime = records[i].EndTime.AddDate(1970, 1, 1)
	}

	return
}
