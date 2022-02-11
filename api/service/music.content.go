package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetMusic ...
func (s Content) GetMusic(id int64) (form form.Music, err error) {
	record := model.Music{ID: id}

	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// ListMusic ...
func (s Content) ListMusic() (records []form.Music, err error) {
	records = []form.Music{}
	err = s.dbc.Limit(10).Order("rank asc").Find(&records).Error

	return
}

// AddMusicFile ...
func (s Content) AddMusicFile(file *multipart.FileHeader) (*model.Upload, error) {

	upload := &model.Upload{}

	src, err := file.Open()
	if err != nil {
		s.log.Error(err)
		return upload, err
	}

	defer src.Close()

	// To get the extension name of the file
	ext := filepath.Ext(file.Filename)

	// To generate a 20 random character generated name for the file and concatenate the extension
	fileName := util.String(20) + ext

	// Destination
	dst, err := os.Create("static/upload/" + fileName)
	if err != nil {
		s.log.Error(err)
		return upload, err
	}

	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		s.log.Error(err)
		return upload, err
	}

	upload.Path = "static/upload/" + fileName
	upload.Size = file.Size

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(upload).Error
	})

	upload.Path = strings.Replace(upload.Path, "static/upload/", "/public/", 1)

	return upload, nil
}

// AddMusicFiles ...
func (s Content) AddMusicFiles(raw json.RawMessage) (path string, err error) {

	upload := &model.Upload{}

	file, err := util.SaveAudioData(s.env, raw)
	if err != nil {
		s.log.Error(err)
		return "", err
	}

	upload.Path = "static/upload/" + file.Name
	upload.Size = file.Size

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(upload).Error
	})

	return upload.Path, nil

}

// AddMusic ...
func (s Content) AddMusic(form *form.Music) (err error) {

	record := &model.Music{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SaveMusic ...
func (s Content) SaveMusic(id int64, form *form.Music) (err error) {

	record := &model.Music{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeleteMusic ...
func (s Content) DeleteMusic(id int64) (err error) {
	record := model.Music{ID: id}

	// Getting the music by id
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	upload := model.Upload{ID: record.UploadID}

	// Getting the upload by id
	if err = s.dbc.Take(&upload).Error; err != nil {
		return
	}

	// Deleting the file
	err = os.Remove(upload.Path)
	if err != nil {
		return err
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(upload).Error
	})

	return
}
