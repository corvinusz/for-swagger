package ctx

import (
	"log"

	"github.com/go-xorm/xorm"
)

// Config is a storage for configuration options
type Config struct {
	Port     string `toml:"port"`
	Database struct {
		Db  string `toml:"db"`
		Dsn string `toml:"dsn"`
	} `toml:"database"`
	Version string
}

// Context is a gate to application services
type Context struct {
	Config Config
	Orm    *xorm.Engine
	Logger *log.Logger
}
