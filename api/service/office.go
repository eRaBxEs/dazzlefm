package service

import (
	"dazzlefm/util"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// Office an interface that satisfies IOffice
type Office struct {
	env *util.Environment
	dbc *gorm.DB
	log *zap.SugaredLogger
}

// compile time guarantee that Office satisfies IOffice
var _ IOffice = Office{}

// NewOffice ...
func NewOffice(env *util.Environment) (service *Office, err error) {
	service = &Office{env: env}
	service.dbc = env.Dbc
	service.log = env.Log.Sugar()

	return
}
