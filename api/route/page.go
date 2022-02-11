package route

// cspell:ignore curri, subscr, signup

import (
	"dazzlefm/service"
	"dazzlefm/service/form"
	"dazzlefm/util"
	"net/http"
	"path"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// Page ...
type Page struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *Page) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("page")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "page")
	grp := env.Rtr.Group(nPre)

	grp.GET("/list", s.GetList)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /page
func (s Page) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, PageModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}
	list, err := s.Content.ListPage()
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// Get GET: /page/id
func (s Page) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Content.GetPage(id)
	if err != nil {
		return
	}

	resp.Set("page", record)
	return c.JSON(200, resp)
}

// Add POST: /page
func (s Page) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Page{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Content.AddPage(form)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /page/id
func (s Page) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Page{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Content.SavePage(id, form)
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /page/id
func (s Page) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Content.DeletePage(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
