package route

// cspell:ignore curri, subscr, signup, jinzhu gorm

import (
	"dazzlefm/service"
	"dazzlefm/service/form"
	"dazzlefm/util"
	"fmt"
	"net/http"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// User ...
type User struct {
	log *zap.SugaredLogger
	env *util.Environment

	Authenticate service.IAuthenticate
}

// Init ...
func (s *User) Init(prefix string, env *util.Environment) error {
	s.env = env
	s.log = env.Log.Sugar().Named("user")
	s.Authenticate = env.Service("auth").(service.IAuthenticate)

	nPre := path.Join(prefix, "user")
	grp := env.Rtr.Group(nPre)

	grp.POST("/login", s.Login)
	grp.GET("/logout", s.Logout)

	grp.GET("/list", s.GetList)
	grp.GET("/:id", s.Get)
	grp.POST("", s.Add)
	grp.POST("/:id", s.Save)
	grp.DELETE("/:id", s.Delete)

	return nil
}

// Login POST: /user/login
func (s User) Login(c echo.Context) (err error) {
	resp := util.Response{Code: 200}

	sMgr, err := util.NewSessionMgr(c, "session")
	if err != nil {
		s.log.Error(err)
		return
	}

	form := &form.Login{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	record, err := s.Authenticate.FindUserByUserName(form.UserName)
	if err == gorm.ErrRecordNotFound {
		resp.Code = http.StatusNotFound
		resp.ErrMsg("username", fmt.Sprintf("invalid user name: %s", form.UserName))
		return c.JSON(200, resp)
	}
	if err != nil {
		s.log.Error(err)
		return
	}

	// password check
	if !util.CheckPasswordHash(form.Password, record.Password) {
		resp.Code = http.StatusNotFound
		resp.ErrMsg("password", fmt.Sprintf("invalid password"))
		return c.JSON(200, resp)
	}

	sMgr.Set("loggedIn", true)
	sMgr.Set("userName", record.UserName)
	sMgr.Set("userID", record.ID)
	sMgr.Set("userRole", record.Role)
	if err = sMgr.Save(); err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("user", record)
	return c.JSON(200, resp)
}

// Logout GET: /user/logout
func (s User) Logout(c echo.Context) (err error) {
	resp := util.Response{Code: 200}

	sMgr, err := util.NewSessionMgr(c, "session")
	if err != nil {
		s.log.Error(err)
		return
	}

	sMgr.Set("loggedIn", false)
	sMgr.Set("userName", "")
	sMgr.Set("userID", int64(0))
	sMgr.Set("userRole", int(0))
	if err = sMgr.Save(); err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("status", "loggedOut")
	return c.JSON(200, resp)
}

// GetList GET: /user
func (s User) GetList(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 1, UserModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: http.StatusOK}

	allowed, err := IsAuthenticated(c)
	if err != nil {
		s.log.Error(err)
		return
	}
	if !allowed {
		resp.Code = http.StatusUnauthorized
		resp.Set("auth", "user not logged in")
		return c.JSON(200, resp)
	}

	list, err := s.Authenticate.ListUser()
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("list", list)
	return c.JSON(200, resp)
}

// Get GET: /user/id
func (s User) Get(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 1, UserModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: http.StatusOK}

	allowed, err := IsAuthenticated(c)
	if err != nil {
		return
	}
	if !allowed {
		resp.Code = http.StatusUnauthorized
		resp.Set("auth", "user not logged in")
		return c.JSON(200, resp)
	}

	id := util.Int64Param(c, "id")

	record, err := s.Authenticate.GetUser(id)
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("user", record)
	return c.JSON(200, resp)
}

// Add POST: /user
func (s User) Add(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 2, UserModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: http.StatusOK}

	allowed, err := IsAuthenticated(c)
	if err != nil {
		return
	}
	if !allowed {
		resp.Code = http.StatusUnauthorized
		resp.Set("auth", "user not logged in")
		return c.JSON(200, resp)
	}

	form := &form.User{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Authenticate.AddUser(form)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.log.Debug(form)

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Save POST: /user/id
func (s User) Save(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 3, UserModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	allowed, err := IsAuthenticated(c)
	if err != nil {
		return
	}
	if !allowed {
		resp.Code = http.StatusUnauthorized
		resp.Set("auth", "user not logged in")
		return c.JSON(200, resp)
	}

	form := &form.User{}
	if err = c.Bind(form); err != nil {
		s.log.Error(err)
		return
	}

	err = s.Authenticate.SaveUser(id, form)
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("id", form.ID)
	return c.JSON(200, resp)
}

// Delete DELETE: /user/id
func (s User) Delete(c echo.Context) (err error) {
	if ok, err := IsAuthorized(c, 4, UserModule); err != nil || !ok {
		if err != nil {
			s.log.Error(err)
			return err
		}
		return c.String(http.StatusNotFound, "!OK")
	}

	resp := util.Response{Code: 200}
	id := util.Int64Param(c, "id")

	allowed, err := IsAuthenticated(c)
	if err != nil {
		return
	}
	if !allowed {
		resp.Code = http.StatusUnauthorized
		resp.Set("auth", "user not logged in")
		return c.JSON(200, resp)
	}

	err = s.Authenticate.DeleteUser(id)
	if err != nil {
		s.log.Error(err)
		return
	}

	resp.Set("id", id)
	return c.JSON(200, resp)
}
