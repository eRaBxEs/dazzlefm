package util

// cSpell:ignore mkr, gocraft, gommon, Sprintf, dbname, Infof
import (
	"fmt"

	"bitbucket.org/mayowa/helpers/dbtools"
	"github.com/jinzhu/gorm"

	"github.com/go-ini/ini"
	"go.uber.org/zap"
)

// InitDb setup the applications database connection
func InitDb(cfg *ini.File, log *zap.SugaredLogger) (db *gorm.DB, err error) {
	log.Info("initializing database connection")

	dbCfg, err := dbtools.GetDbConfig("db", cfg)
	if err != nil {
		log.Error(err)
		return
	}

	db, err = gorm.Open(dbCfg.Driver, dbCfg.DSN())
	if err != nil {
		log.Error(err)
		return
	}

	if db == nil {
		err = fmt.Errorf("nil return: dbr.Open(%s, %s)", dbCfg.Driver, dbCfg.String())
		log.Error(err)
		return
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)
	db.SingularTable(true)

	err = db.DB().Ping()
	if err != nil {
		log.Error(err)
		log.Debugf("dbr.Open(%s, %s)", dbCfg.Driver, dbCfg.DSN())
		return
	}

	log.Infof("connected to db with: driver:%s, dbname:%s", dbCfg.Driver, dbCfg.MakeString(true))

	return
}

type txFunc func(*gorm.DB) error

// Transact is a closure that wraps a transaction
func Transact(db *gorm.DB, log *zap.SugaredLogger, fn txFunc) error {
	tx := db.Begin()
	if err := db.Error; err != nil {
		log.Error(err)
		return err
	}

	if err := fn(tx); err != nil {
		log.Error(err)

		tx.Rollback()
		if err := tx.Error; err != nil {
			log.Error(err)
			return err
		}

		return err
	}

	tx.Commit()
	if err := tx.Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
