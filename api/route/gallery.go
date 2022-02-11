package route

// cspell:ignore curri, subscr, signup

import (
	"dazzlefm/service"
	"dazzlefm/service/form"
	"dazzlefm/util"
	"net/http"
	"path"
	"strconv"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// Gallery ...
type Gallery struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *Gallery) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("gallery")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "gallery")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /gallery/
func (s Gallery) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, GalleryModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	page := c.QueryParam("page")

	pageNo, err := strconv.Atoi(page)
	if err != nil {
		return
	}

	offSet := c.QueryParam("offset")

	offset, err := strconv.Atoi(offSet)
	if err != nil {
		return
	}

	resp := util.Response{Code: 200}
	list, err := s.Content.ListGallery(pageNo, offset)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount GET: /gallery/count
func (s Gallery) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Content.CountGallery()
	if err != nil {
		return
	}

	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// Get GET: /gallery/:id
func (s Gallery) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Content.GetGallery(id)
	if err != nil {
		return
	}

	resp.Set("gallery", record)
	return c.JSON(200, resp)
}

// Add POST: /gallery
func (s Gallery) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, GalleryModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Gallery{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.AddGallery(form)
	if err != nil {
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /gallery/:id
func (s Gallery) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, GalleryModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Gallery{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.SaveGallery(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /gallery/:id
func (s Gallery) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, GalleryModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Content.DeleteGallery(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
