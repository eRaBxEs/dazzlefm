package util

import (
	"strings"

	"bitbucket.org/mayowa/helpers/validate"
	"gopkg.in/go-playground/validator.v9"
)

// Response A map[string]interface{} type for responding to json clients
type Response struct {
	Code   int                    `json:"code"`
	Store  map[string]interface{} `json:"store"`
	Errors ErrMsg                 `json:"errors"`
}

// ErrMsg holds error messages
type ErrMsg map[string]string

// Set stores an a named value in Data
func (s *Response) Set(name string, value interface{}) {
	if s.Store == nil {
		s.Store = map[string]interface{}{}
	}

	s.Store[name] = value
}

// ErrMsg helper for adding error messages
func (s *Response) ErrMsg(name, value string) {
	// create error field if it doesn't exist

	if s.Errors == nil {
		s.Errors = ErrMsg{}
	}

	s.Errors[name] = value
}

// ValidationErrors process validator.ValidationErrors
func (s *Response) ValidationErrors(err error) {
	ver, ok1 := err.(validator.ValidationErrors)
	ler, ok2 := err.(ValidationError)
	if !ok1 && !ok2 {
		s.ErrMsg("error", err.Error())
		return
	}

	if ok1 {
		for _, e := range ver {
			p := strings.Split(e.Field(), ".")
			fNme := p[len(p)-1]
			s.ErrMsg(fNme, validate.MakeErrMessage(e))
		}
	} else if ok2 {
		s.ErrMsg(ler.Field, ler.ErrMsg)
	}
}

// DbErrors attempts to detect a db error
func (s *Response) DbErrors(err error) bool {
	errMsg := err.Error()

	if strings.Contains(errMsg, "UNIQUE constraint failed") {
		s.ErrMsg("error", errMsg)
		return true
	}

	return false
}
