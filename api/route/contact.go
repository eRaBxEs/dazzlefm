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

// Contact ...
type Contact struct {
	log *zap.SugaredLogger
	env *util.Environment

	Office service.IOffice
}

// Init ...
func (s *Contact) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("contact")
	s.Office = env.Service("office").(service.IOffice)

	nPre := path.Join(prefix, "contact")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /Contact
func (s Contact) GetList(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 1, ContactsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

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
	list, err := s.Office.ListContact(pageNo, offset)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount GET: /contact/count
func (s Contact) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Office.CountContact()
	if err != nil {
		return
	}

	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// Get GET: /Contact/id
func (s Contact) Get(c echo.Context) (err error) {

	if ok, err := IsAuthorized(c, 1, ContactsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Office.GetContact(id)
	if err != nil {
		return
	}

	resp.Set("contact", record)
	return c.JSON(200, resp)
}

// Add POST: /Contact
func (s Contact) Add(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 2, ContactsModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}

	form := &form.Contact{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Office.AddContact(form)
	if err != nil {
		return
	}

	//   s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /Contact/id
func (s Contact) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, ContactsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Contact{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Office.SaveContact(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /Contact/id
func (s Contact) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, ContactsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Office.DeleteContact(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
