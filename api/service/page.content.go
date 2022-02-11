package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"dazzlefm/util"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

// GetPage ...
func (s Content) GetPage(id int64) (form form.Page, err error) {
	record := model.Page{ID: id}
	if err = s.dbc.Take(&record).Error; err != nil {
		return
	}

	err = copier.Copy(&form, record)

	return
}

// ListPage ...
func (s Content) ListPage() (records []form.Page, err error) {
	records = []form.Page{}
	err = s.dbc.Order("name asc").Find(&records).Error

	return
}

// AddPage ...
func (s Content) AddPage(form *form.Page) (err error) {

	form.Attributes, err = s.processCarouselImage(form.Attributes)
	if err != nil {
		return err
	}

	record := &model.Page{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Create(record).Error
	})

	form.ID = record.ID

	return
}

// SavePage ...
func (s Content) SavePage(id int64, form *form.Page) (err error) {

	form.Attributes, err = s.processCarouselImage(form.Attributes)
	if err != nil {
		return err
	}

	record := &model.Page{}
	if err = copier.Copy(record, form); err != nil {
		return
	}

	record.ID = id
	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Save(record).Error
	})

	return
}

// DeletePage ...
func (s Content) DeletePage(id int64) (err error) {
	record := model.Page{ID: id}

	err = util.Transact(s.dbc, s.log, func(tx *gorm.DB) error {
		return tx.Delete(record).Error
	})

	return
}

func (s Content) processImage(attributes json.RawMessage) (data []byte, err error) {
	attribs := map[string]interface{}{}
	if err = json.Unmarshal(attributes, &attribs); err != nil {
		s.log.Error(err)
		return
	}

	images, _ := attribs["images"].(map[string]interface{})

	// s.log.Debugf("images:%v", images)

	for _, i := range images {
		imgData, err := json.Marshal(i)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		s.log.Debugf("imgData: %v", imgData)

		image, err := util.SaveImageData(s.env, imgData)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		s.log.Debugf("%d", image.Size)

		i.(map[string]interface{})["data"] = image.Data
		i.(map[string]interface{})["name"] = image.Name
		i.(map[string]interface{})["size"] = image.Size
	}

	s.log.Debugf("attribs:%v", attribs)
	return json.Marshal(attribs)

}

func (s Content) processCarouselImage(attributes json.RawMessage) (data []byte, err error) {
	// s.log.Debugf("record:%v", attributes)

	attribs := map[string]interface{}{}
	if err = json.Unmarshal(attributes, &attribs); err != nil {
		s.log.Error(err)
		return
	}

	carouseler := attribs["carousel"]

	// if carousel == nil {
	// 	s.log.Debugf("\ncarousel is very nil\n")
	// }

	if carouseler != nil {
		s.log.Debugf("\ncarousel not nil\n")

		carousel := attribs["carousel"].([]interface{})

		for _, i := range carousel {
			realImage := i.(map[string]interface{})["image"]
			imgData, err := json.Marshal(realImage)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}
			s.log.Debugf("imgData: %v", imgData)

			image, err := util.SaveImageData(s.env, imgData)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}
			s.log.Debugf("size:%d", image.Size)

			realImage.(map[string]interface{})["data"] = image.Data
			realImage.(map[string]interface{})["name"] = image.Name
			realImage.(map[string]interface{})["size"] = image.Size
		}
	}

	lead := attribs["leaders"]

	if lead != nil {

		leaders := attribs["leaders"].([]interface{})

		s.log.Debugf("\nleaders is not nil\n")

		for _, i := range leaders {

			realImage := i.(map[string]interface{})["image"]
			imgData, err := json.Marshal(realImage)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}

			// s.log.Debugf("imgData: %v", imgData)

			image, err := util.SaveImageData(s.env, imgData)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}

			s.log.Debugf("size:%d", image.Size)

			realImage.(map[string]interface{})["data"] = image.Data
			realImage.(map[string]interface{})["name"] = image.Name
			realImage.(map[string]interface{})["size"] = image.Size
		}

	}

	//s.log.Debugf("attribs:%v", attribs)
	return json.Marshal(attribs)
}

func (s Content) processLeaderImage(attributes json.RawMessage) (data []byte, err error) {
	// s.log.Debugf("record:%v", attributes)

	attribs := map[string]interface{}{}
	if err = json.Unmarshal(attributes, &attribs); err != nil {
		s.log.Error(err)
		return
	}

	leaders := attribs["leaders"].([]interface{})
	if leaders != nil {

		s.log.Debugf("\nleaders is not nil\n")

		for _, i := range leaders {

			realImage := i.(map[string]interface{})["image"]
			imgData, err := json.Marshal(realImage)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}

			// s.log.Debugf("imgData: %v", imgData)

			image, err := util.SaveImageData(s.env, imgData)
			if err != nil {
				s.log.Error(err)
				return nil, err
			}

			s.log.Debugf("size:%d", image.Size)

			realImage.(map[string]interface{})["data"] = image.Data
			realImage.(map[string]interface{})["name"] = image.Name
			realImage.(map[string]interface{})["size"] = image.Size
		}

	}

	if leaders == nil {
		s.log.Debugf("\nleaders is very nil\n")
	}

	s.log.Debugf("attribs:%v", attribs)
	return json.Marshal(attribs)
}
