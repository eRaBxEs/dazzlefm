package util

import (
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// IHandler specs handlers used by mkr
type IHandler interface {
	Init(string, *Environment) error
}

// Environment variables that represent the executing environment
type Environment struct {
	AppName     string
	Rtr         *echo.Echo
	Dbc         *gorm.DB
	Cfg         *ini.File
	Log         *zap.Logger
	Validate    *validator.Validate
	Paths       map[string]string
	Urls        map[string]string
	ServiceList map[string]interface{}
}

// Path returns a path from Environment.Paths or "" if not found
func (s Environment) Path(name string) string {
	val, _ := s.Paths[name]

	return val
}

// URL returns a url from Environment.Urls or "" if not found
func (s Environment) URL(name string) string {
	val, _ := s.Urls[name]

	return val
}

// FullURL returns a url with public_base_url prefixed from Environment.Urls or "" if not found
func (s Environment) FullURL(name string) string {
	val, _ := s.Urls[name]

	if len(val) > 0 {
		val = s.Cfg.Section("url").Key("base_url").String() + val
	}

	return val
}

// PublicURL returns a url with public_base_url prefixed from Environment.Urls or "" if not found
func (s Environment) PublicURL(name string) string {
	val, _ := s.Urls[name]

	if len(val) > 0 {
		val = s.Cfg.Section("url").Key("public_base_url").String() + val
	}

	return val
}

// Service ...
func (s Environment) Service(name string) interface{} {
	_, exists := s.ServiceList[name]
	if !exists {
		return nil
	}

	return s.ServiceList[name]
}
