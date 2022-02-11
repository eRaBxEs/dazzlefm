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

// Script ...
type Script struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *Script) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("script")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "snippets/script")
	grp := env.Rtr.Group(nPre)

	grp.GET("/", s.GetList)
	grp.GET("/:id", s.Get)
	grp.POST("/", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /snippets/script
func (s Script) GetList(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	list, err := s.Content.ListSnippet()
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// Get GET: /snippets/script/id
func (s Script) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Content.GetSnippet(id)
	if err != nil {
		return
	}

	resp.Set("script", record)
	return c.JSON(200, resp)
}

// Add POST: /snippets/script
func (s Script) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, ScriptModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Script{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.AddSnippet(form)
	if err != nil {
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /snippets/script/id
func (s Script) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, ScriptModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Script{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.SaveSnippet(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /snippets/script/id
func (s Script) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, ScriptModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Content.DeleteSnippet(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
