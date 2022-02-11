package service

import (
	"dazzlefm/util"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// Authenticate an interface that satisfies IAuthenticate
type Authenticate struct {
	env *util.Environment
	dbc *gorm.DB
	log *zap.SugaredLogger
}

// compile time guarantee that Authenticate satisfies IAuthenticate
var _ IAuthenticate = Authenticate{}

// NewAuthenticate ...
func NewAuthenticate(env *util.Environment) (service *Authenticate, err error) {
	service = &Authenticate{env: env}
	service.dbc = env.Dbc
	service.log = env.Log.Sugar()

	return
}
