package util

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var gValidate *validator.Validate

// Validator an interface for models that can be validated by the Validate function
type Validator interface {
	Validate(*gorm.DB) error
}

// ValidationError the error type returned by Validater
type ValidationError struct {
	Field  string
	ErrMsg string
}

// Error implements the Error interface
func (s ValidationError) Error() string {
	return s.ErrMsg
}

// Validate binds and validates a struct
func Validate(env *Environment, c echo.Context, tx *gorm.DB, obj interface{}) (err error) {

	if err = c.Bind(obj); err != nil {
		return
	}

	if err = ValidateOnly(env, tx, obj); err != nil {
		return
	}

	return nil
}

// ValidateOnly ...
func ValidateOnly(env *Environment, tx *gorm.DB, obj interface{}) (err error) {

	err = env.Validate.Struct(obj)
	if err != nil {
		return
	}

	vObj, ok := obj.(Validator)
	if ok {
		err = vObj.Validate(tx)
		return
	}

	return nil
}
