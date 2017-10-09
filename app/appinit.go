package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-xorm/xorm"
	"github.com/pelletier/go-toml"

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
	a.C.Logger = log.New(os.Stdout, os.Args[0]+":", log.Flags())

	// init ORM
	return a.initORM()
}

//------------------------------------------------------------------------------
// initConfigFromFile reads configuration file and also sets default parameters
func (a *App) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return fmt.Errorf("configuration file %s read error: %s", cfgFileName, err.Error())
	}
	err = toml.Unmarshal(tomlData, &a.C.Config)
	if err != nil {
		return fmt.Errorf("configuration file %s decoding error: %s", cfgFileName, err.Error())
	}

	// set config default values
	if len(a.C.Config.Port) == 0 {
		a.C.Config.Port = "11011"
	}
	// set version
	a.C.Config.Version = "0.0.1"
	return nil
}

//------------------------------------------------------------------------------
// init database from known application config
func (a *App) initORM() error {
	var err error
	// create
	a.C.Orm, err = xorm.NewEngine(a.C.Config.Database.Db, a.C.Config.Database.Dsn)
	if err != nil {
		return fmt.Errorf("cannot connect to db %s, error: %s", a.C.Config.Database.Dsn, err.Error())
	}
	// SQL log on
	a.C.Orm.ShowExecTime(true)
	a.C.Orm.ShowSQL(true)
	// sync schema
	return a.syncSchema()
}

//------------------------------------------------------------------------------
// sync Schema
func (a *App) syncSchema() error {
	// migrate tables
	// if err = a.Context.Orm.Sync(new(DbConstant)); err != nil {
	// 	return err
	// }
	err := a.C.Orm.Sync(new(groups.Entity))
	if err != nil {
		return fmt.Errorf("cannot operate table 'groups', error: %s", err.Error())
	}
	err = a.C.Orm.Sync(new(users.Entity))
	if err != nil {
		return fmt.Errorf("cannot operate table 'users', error: %s", err.Error())
	}

	return nil
}
