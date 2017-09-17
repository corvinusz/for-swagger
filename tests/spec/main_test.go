package spectests

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/go-resty/resty"
	"github.com/go-testfixtures/testfixtures"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sclevine/spec"

	"github.com/corvinusz/for-swagger/app"
)

func TestAppStart(t *testing.T) {
	spec.Run(t, "appStart", func(t *testing.T, when spec.G, it spec.S) {
		testAppStart(t, when, it)
	}, spec.Flat())
}

func testAppStart(t *testing.T, when spec.G, it spec.S) {
	var (
		suite *AppTestSuite
		err   error
	)

	it.Before(func() {
		suite = new(AppTestSuite)
		err = suite.setup()
		if err != nil {
			t.Error(err.Error())
		}
	})

	it("should respond to dial", func() {
		err = suite.waitServerStart(3 * time.Second)
		if err != nil {
			t.Error(err.Error())
		}
	})
}

// AppTestSuite is testing context for suite
type AppTestSuite struct {
	app     *app.App
	baseURL string
	client  *resty.Client
}

var suite *AppTestSuite

const (
	// cfgFileName    = "../test-config/test-config.toml" // this is correct line
	cfgFileName    = "./test-config/test-config.toml" // this is incorrect line
	fixturesFolder = "../fixtures"
)

// setup called once before test
func (s *AppTestSuite) setup() error {
	err := s.setupServer()
	if err != nil {
		return err
	}
	s.baseURL = "http://localhost:" + s.app.Context.Config.Port
	// create and setup resty client
	s.client = resty.DefaultClient
	s.client.SetHeader("Content-Type", "application/json")
	s.client.SetHostURL(s.baseURL)
	return nil
}

//------------------------------------------------------------------------------
func (s *AppTestSuite) setupServer() error {
	var err error
	// init test application
	s.app, err = app.New(app.CmdFlags{CfgFileName: cfgFileName})
	if err != nil {
		return err
	}
	// load fixtures
	err = testfixtures.LoadFixtures(fixturesFolder, s.app.Context.Orm.DB().DB, &testfixtures.SQLite{})
	if err != nil {
		return err
	}
	// start test server with go routine
	go s.app.Run()
	// return
	return nil
	// wait til server started then return
	//s.waitServerStart(3 * time.Second)
}

//------------------------------------------------------------------------------
func (s *AppTestSuite) waitServerStart(timeout time.Duration) error {
	done := time.Now().Add(timeout)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", ":"+s.app.Context.Config.Port)
		if err == nil {
			return c.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("cannot connect %v for %v", s.baseURL, timeout)
}
