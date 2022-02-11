package util

import (
	"fmt"
	"os"

	"bitbucket.org/mayowa/helpers/config"
	"bitbucket.org/mayowa/helpers/path"

	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// InitLogger ...
func InitLogger() *zap.Logger {

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}
	logger, err := cfg.Build()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	zap.ReplaceGlobals(logger)

	return logger
}

// CreateDefaultConfig ...
func CreateDefaultConfig(name, cfgPath string, log *zap.SugaredLogger) (cfg *ini.File, err error) {
	log.Info("creating default confing in: ", cfgPath)

	cfgStr := fmt.Sprintf(`
		cfgpath = %s
		host =
		port = 4000

		[url]
		origin = "http://localhost:3000", "http://localhost:4000"
		base_url =
		public_base_url = "http://localhost:3000"
		static = /static
		upload = /public
		downloads = /downloads

		[path]
		static = static
		upload = static/upload
		downloads = static/downloads

		[db]
		driver   = postgres
		host     = localhost:5432
		user     = sysdba
		password = dev@dba@4839
		dbname   = dazzlefmdb
		sslmode  = disable
		`,
		cfgPath,
	)

	if cfg, _, err = config.Create(name, cfgStr); err != nil {
		log.Error(err)
		return
	}

	return
}

// InitConfig load config file and create defaults if config not found
func InitConfig(appName string, log *zap.SugaredLogger) (cfg *ini.File, cfgPath string, err error) {
	cfgPath, err = config.GetPath(appName)
	if err != nil {
		return
	}
	// found a config location
	cfg, cfgPath, err = config.Load(appName)
	if err != nil {
		if len(cfgPath) == 0 {
			log.Error(err)
			return
		}

		if cfg, err = CreateDefaultConfig(appName, cfgPath, log); err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Info("loading config from: ", cfgPath)
	}

	log.Infof("CORS allowed-origins: %s", cfg.Section("url").Key("origin").Strings(","))
	log.Infof("base_url: %s", cfg.Section("url").Key("base_url").String())

	// check if paths exist
	paths := cfg.Section("path").KeysHash()
	for k, pth := range paths {
		if !path.Available(pth) {
			log.Errorf("path[%s]: %s not accessible", k, pth)
			if err := os.MkdirAll(pth, 0700); err != nil {
				log.Errorf("path[%s]: %s cannot create folder", k, pth)
			}
		} else {
			log.Infof("path[%s]: %s", k, pth)
		}
	}

	return
}

// MakeEnv a container for application globals to be passed around within the app
func MakeEnv(appName string, rtr *echo.Echo, cfg *ini.File, dbc *gorm.DB, lgr *zap.Logger) (env *Environment) {
	env = &Environment{
		AppName:     appName,
		Rtr:         rtr,
		Dbc:         dbc,
		Cfg:         cfg,
		Log:         lgr,
		Validate:    validator.New(),
		Paths:       cfg.Section("path").KeysHash(),
		Urls:        cfg.Section("url").KeysHash(),
		ServiceList: map[string]interface{}{},
	}

	return
}
