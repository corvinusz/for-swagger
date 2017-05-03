package app

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-xorm/xorm"

	"github.com/corvinusz/for-swagger/server/groups"
	"github.com/corvinusz/for-swagger/server/users"
)

//------------------------------------------------------------------------------
func (a *App) initApp(f CmdFlags) error {
	// read config file
	err := a.initConfigFromFile(f.CfgFileName)
	if err != nil {
		return err
	}

	// init Logger
	a.Context.Logger = log.New(os.Stdout, "RERS: ", log.Flags())

	// init ORM
	return a.initORM()
}

//------------------------------------------------------------------------------
// initConfigFromFile reads configuration file and also sets default parameters
func (a *App) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return errors.New("Configuration file read error: " + cfgFileName + "\nError:" + err.Error())
	}
	_, err = toml.Decode(string(tomlData[:]), &a.Context.Config)
	if err != nil {
		return errors.New("Configuration file decoding error: " + cfgFileName + "\nError:" + err.Error())
	}

	// set config default values
	if len(a.Context.Config.Port) == 0 {
		a.Context.Config.Port = "11011"
	}
	// set version
	a.Context.Config.Version = "0.0.1"
	return nil
}

//------------------------------------------------------------------------------
// init database from known application config
func (a *App) initORM() error {
	var err error
	// create
	a.Context.Orm, err = xorm.NewEngine(a.Context.Config.Database.Db, a.Context.Config.Database.Dsn)
	if err != nil {
		return err
	}
	// SQL on
	a.Context.Orm.ShowExecTime(true)
	a.Context.Orm.ShowSQL(true)
	// sync schema
	return a.syncSchema()
}

//------------------------------------------------------------------------------
// sync Schema
func (a *App) syncSchema() error {
	var err error
	// migrate tables
	// if err = a.Context.Orm.Sync(new(DbConstant)); err != nil {
	// 	return err
	// }
	if err = a.Context.Orm.Sync(new(groups.Entity)); err != nil {
		return err
	}
	if err = a.Context.Orm.Sync(new(users.Entity)); err != nil {
		return err
	}

	return err
}
