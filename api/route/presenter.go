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

// Presenter ...
// NOTE: The miracle that happens here in simple words is below:
// Embedding an interface as an attribute of a type that satisfies a different interface
type Presenter struct {
	log *zap.SugaredLogger
	env *util.Environment

	Office service.IOffice
}

// Init ...
func (s *Presenter) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("presenter")
	s.Office = env.Service("office").(service.IOffice)
	// Here I used type assertion to set the type of the empty interface (interface{})

	nPre := path.Join(prefix, "presenter")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /presenter
func (s Presenter) GetList(c echo.Context) (err error) {
	// if ok, err := IsAuthorized(c, 1, PresenterModule); err != nil || !ok {
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

	statusString := c.QueryParam("status")

	status, err := strconv.Atoi(statusString)
	if err != nil {
		return
	}

	resp := util.Response{Code: 200}
	list, err := s.Office.ListPresenter(pageNo, offset, status)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount GET: /presenter/count
func (s Presenter) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Office.CountPresenter()
	if err != nil {
		return
	}

	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// Get GET: /presenter/id
func (s Presenter) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Office.GetPresenter(id)
	if err != nil {
		return
	}

	resp.Set("presenter", record)
	return c.JSON(200, resp)
}

// Add POST: /presenter
func (s Presenter) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, PresenterModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Presenter{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Office.AddPresenter(form)
	if err != nil {
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /presenter/id
func (s Presenter) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, PresenterModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Presenter{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Office.SavePresenter(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /presenter/id
func (s Presenter) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, PresenterModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Office.DeletePresenter(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
