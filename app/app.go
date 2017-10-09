package app

import (
	"github.com/corvinusz/for-swagger/ctx"
	"github.com/corvinusz/for-swagger/server"
)

// App is an application with context
type App struct {
	C   *ctx.Context
	Srv *server.Server
}

// CmdFlags represents command run flags
type CmdFlags struct {
	CfgFileName string
}

// New is a constructor
func New(f CmdFlags) (*App, error) {
	// construct
	a := new(App)
	a.C = new(ctx.Context)
	// init
	err := a.initApp(f)
	if err != nil {
		return nil, err
	}
	// report
	a.C.Logger.Printf("successfully initialized: %+v\n", a.C.Config)
	return a, nil
}

// Run starts an application
func (a App) Run() {
	// create server and run it
	a.Srv = server.New(a.C)
	a.Srv.Start()
}

// Shutdown gracefully stops server
func (a App) Shutdown() {
	// stop server
	a.C.Logger.Println("appcontrol", "stopping server")
	if a.Srv != nil {
		a.Srv.Shutdown()
	}

	// close database connection
	if a.C.Orm != nil {
		a.C.Logger.Println("appcontrol", "closing db connection")
		a.C.Orm.Close()
	}
}
