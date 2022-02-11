//
// session.go
// helper functions to simplify session management
// Copyright 2017 Akinmayowa Akinyemi
//

package util

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo"
)

// SessionMgr utillity to manage sessions
type SessionMgr struct {
	session *sessions.Session
	ctx     echo.Context
}

// NewSessionMgr create an instance of SessionMgr
func NewSessionMgr(c echo.Context, sessName string) (sMgr *SessionMgr, err error) {
	sMgr = new(SessionMgr)
	sMgr.session, _ = session.Get(sessName, c)
	if err != nil {
		return nil, err
	}

	sMgr.ctx = c

	return
}

// Set store a value into the session
func (s *SessionMgr) Set(key string, value interface{}) {
	s.session.Values[key] = value
}

// Value gets an item from the session. performs check for existence
func (s SessionMgr) Value(key string) (val interface{}, exists bool) {

	if _, exists = s.session.Values[key]; !exists {
		return nil, false
	}

	return s.session.Values[key], true
}

// Save ...
func (s *SessionMgr) Save() error {
	if err := s.session.Save(s.ctx.Request(), s.ctx.Response()); err != nil {
		return err
	}

	return nil
}

// StringValue string typed version of Value
func (s SessionMgr) StringValue(key string) (val string, exists bool) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(string), true
}

// String ...
func (s SessionMgr) String(key string) (val string) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(string)
}

// IntValue int typed version of Value
func (s SessionMgr) IntValue(key string) (val int, exists bool) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(int), true
}

// Int ...
func (s SessionMgr) Int(key string) (val int) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(int)
}

// Int64Value an int64 typed version of Value
func (s SessionMgr) Int64Value(key string) (val int64, exists bool) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(int64), true
}

// Int64 ...
func (s SessionMgr) Int64(key string) (val int64) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(int64)
}

// BoolValue an bool typed version of Value
func (s SessionMgr) BoolValue(key string) (val bool, exists bool) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(bool), true
}

// Bool ...
func (s SessionMgr) Bool(key string) (val bool) {
	retv, exists := s.Value(key)
	if !exists {
		return
	}

	return retv.(bool)
}
