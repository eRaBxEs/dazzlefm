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

// Event ...
type Event struct {
	log *zap.SugaredLogger
	env *util.Environment

	Office service.IOffice
}

// Init ...
func (s *Event) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("event")
	s.Office = env.Service("office").(service.IOffice)

	nPre := path.Join(prefix, "event")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/count", s.GetListCount)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	grp.GET("/live", s.GetLive)
	grp.GET("/others", s.GetOthers)

	return nil
}

// GetList GET: /events
func (s Event) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, EventModule); err != nil || !ok {
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
	list, err := s.Office.ListEvent(pageNo, offset)
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// GetListCount GET: /event/count
func (s Event) GetListCount(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	cnt, err := s.Office.CountEvent()
	if err != nil {
		return
	}

	resp.Set("count", cnt)
	return c.JSON(200, resp)
}

// Get GET: /events/id
func (s Event) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Office.GetEvent(id)
	if err != nil {
		return
	}

	resp.Set("event", record)
	return c.JSON(200, resp)
}

// Add POST: /event
func (s Event) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, EventModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Event{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Office.AddEvent(form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /event/id
func (s Event) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, EventModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Event{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Office.SaveEvent(id, form)
	if err != nil {
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /event/id
func (s Event) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, EventModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Office.DeleteEvent(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}

// GetLive GET: /event/live
func (s Event) GetLive(c echo.Context) (err error) {
	// if ok, err := IsAuthorized(c, 1, EventModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}

	eventObj, err := s.Office.GetLiveEvent()
	if err != nil {
		return
	}

	resp.Set("live", eventObj)
	return c.JSON(200, resp)
}

// GetOthers GET: /event/others
func (s Event) GetOthers(c echo.Context) (err error) {
	// if ok, err := IsAuthorized(c, 1, EventModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}

	events, err := s.Office.GetOtherEvents()
	if err != nil {
		return
	}

	resp.Set("others", events)
	return c.JSON(200, resp)
}
