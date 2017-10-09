package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"

	_ "github.com/mattn/go-sqlite3"

	"github.com/corvinusz/for-swagger/app"
)

var (
	gopath     = os.Getenv("GOPATH")
	configFlag = flag.String("config",
		gopath+"/src/github.com/corvinusz/for-swagger/config/config.toml",
		"-config=\"path-to-your-config-file\" ")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// parse flags
	flag.Parse()
	flags := app.CmdFlags{
		CfgFileName: *configFlag,
	}
	// create app
	a, err := app.New(flags)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// here we go
		a.Run()
	}()

	// signal control
	sigstop := make(chan os.Signal)
	signal.Notify(sigstop, os.Kill, os.Interrupt)

	sig := <-sigstop // wait while server works
	// caught stop signal
	if a.C.Logger != nil {
		a.C.Logger.Println("appcontrol", os.Args[0]+" caught signal "+sig.String())
	}

	// shutdown server
	a.Shutdown()
}
