package service

import (
	"dazzlefm/util"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// Content an interface that satisfies IContent
type Content struct {
	env *util.Environment
	dbc *gorm.DB
	log *zap.SugaredLogger
}

// compile time guarantee that Content satisfies IContent
var _ IContent = Content{}

// NewContent ...
func NewContent(env *util.Environment) (service *Content, err error) {
	service = &Content{env: env}
	service.dbc = env.Dbc
	service.log = env.Log.Sugar()

	return
}
