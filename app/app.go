package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/corvinusz/for-swagger/ctx"
	"github.com/corvinusz/for-swagger/server"
)

// App is an application with context
type App struct {
	Context *ctx.Context
	Srv     *server.Server
}

// CmdFlags represents command run flags
type CmdFlags struct {
	CfgFileName string
}

// New is a constructor
func New(f CmdFlags) (*App, error) {
	// construct
	a := new(App)
	a.Context = new(ctx.Context)
	// init
	err := a.initApp(f)
	if err != nil {
		return nil, err
	}
	// report
	a.Context.Logger.Printf("successfully initialized: %+v\n", a.Context.Config)
	return a, nil
}

// Run starts an application
func (a *App) Run() {
	// set up os-signal catchers channel
	var fail bool
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// start os-signal catch goroutine
	go func() {
		for sig := range sigChan {
			if a.Srv != nil {
				err := a.Srv.Shutdown()
				if err != nil {
					fail = true
					log.Println(os.Args[0] + " shutdown error:\n " + err.Error())
				}
			}
			if a.Context.Orm.DB() != nil {
				err := a.Context.Orm.Close()
				if err != nil {
					fail = true
					log.Println(os.Args[0] + " shutdown error:\n " + err.Error())
				}
			}
			if fail {
				os.Exit(1)
			}
			log.Println(os.Args[0] + " gracefully stopped on " + sig.String())
			os.Exit(0)
		}
	}()
	// create server and run it
	a.Srv = server.New(a.Context)
	a.Srv.Start()
}
