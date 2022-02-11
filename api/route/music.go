package route

// cspell:ignore curri, subscr, signup

import (
	"dazzlefm/service"
	"dazzlefm/service/form"
	"dazzlefm/util"
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// Music ...
type Music struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *Music) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("music")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "music")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	grp.GET("/:id", s.Get)
	grp.POST("/upload", s.SaveFile)
	grp.POST("/uploads", s.SaveFiles)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /music
func (s Music) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, PageModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}
	list, err := s.Content.ListMusic()
	if err != nil {
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// Get GET: /music/id
func (s Music) Get(c echo.Context) (err error) {

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	record, err := s.Content.GetMusic(id)
	if err != nil {
		return
	}

	resp.Set("music", record)
	return c.JSON(200, resp)
}

// SaveFile POST: /music/upload
func (s Music) SaveFile(c echo.Context) (err error) {

	if ok, err := IsAuthorized(c, 2, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	log.Printf("%#v", file.Header)

	res, err := s.Content.AddMusicFile(file)
	if err != nil {
		return
	}

	resp.Set("upload", res)
	return c.JSON(200, resp)
}

// SaveFiles POST: /music/uploads (Not working yet ...)
func (s Music) SaveFiles(c echo.Context) (err error) {

	if ok, err := IsAuthorized(c, 2, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	raw := &util.ImageData{}

	if err = c.Bind(raw); err != nil {
		return
	}

	rawFile, err := json.Marshal(raw)
	if err != nil {
		return
	}

	res, err := s.Content.AddMusicFiles(rawFile)
	if err != nil {
		return
	}

	resp.Set("path", res)
	return c.JSON(200, resp)

}

// Add POST: /music
func (s Music) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Music{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Content.AddMusic(form)
	if err != nil {
		s.log.Error(err)
		return
	}

	// s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /music/id
func (s Music) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	form := &form.Music{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Content.SaveMusic(id, form)
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /music/id
func (s Music) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, PageModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	err = s.Content.DeleteMusic(id)
	if err != nil {
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
