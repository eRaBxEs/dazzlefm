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

// News ...
type News struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *News) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("news")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "news")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /news
func (s News) GetList(c echo.Context) (err error) {
	// if ok, err := IsAuthorized(c, 1, NewsModule); err != nil || !ok {
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
	list, err := s.Content.ListNews(pageNo, offset)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount ...
func (s News) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Content.CountNews()
	if err != nil {
		return
	}
	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// Get GET: /news/id
func (s News) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Content.GetNews(id)
	if err != nil {
		return
	}

	resp.Set("news", record)
	return c.JSON(200, resp)
}

// Add POST: /news
func (s News) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, NewsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.News{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.AddNews(form)
	if err != nil {
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /news/id
func (s News) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, NewsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.News{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.SaveNews(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /news/id
func (s News) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, NewsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Content.DeleteNews(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
