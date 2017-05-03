package main

import (
	"flag"
	"log"
	"os"
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
	// run
	a.Run()
}
