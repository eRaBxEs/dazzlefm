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

// Settings ...
type Settings struct {
	log *zap.SugaredLogger
	env *util.Environment

	Content service.IContent
}

// Init ...
func (s *Settings) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("settings")
	s.Content = env.Service("content").(service.IContent)

	nPre := path.Join(prefix, "settings")
	grp := env.Rtr.Group(nPre)

	grp.GET("", s.GetList)
	// grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	// grp.POST("/:id", s.Save)
	// grp.DELETE("/:id", s.Delete)

	return nil
}

// GetList GET: /settings
func (s Settings) GetList(c echo.Context) (err error) {

	// if ok, err := IsAuthorized(c, 1, SettingsModule); err != nil || !ok {
	// 	if err != nil {
	// 		s.log.Error(err)
	// 		return err
	// 	}
	// 	return c.String(http.StatusNotFound, "!OK")
	// }

	resp := util.Response{Code: 200}
	list, err := s.Content.ListSettings()
	if err != nil {
		return
	}

	resp.Set("settings", list)
	return c.JSON(200, resp)
}

// Add POST: /settings
func (s Settings) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, SettingsModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}

	form := &form.Settings{}
	if err = c.Bind(form); err != nil {
		return
	}

	err = s.Content.AddSettings(form)
	if err != nil {
		return
	}

	s.log.Debug(form)

	//resp.Set("id", form.ID)
	return c.JSON(200, resp)
}
