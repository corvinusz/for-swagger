package bddtests

import (
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-resty/resty"
	"github.com/go-testfixtures/testfixtures"
	_ "github.com/mattn/go-sqlite3"

	"github.com/corvinusz/for-swagger/app"
)

func TestBddsuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bddtests Suite")
}

// RersTestSuite is testing context for suite
type RersTestSuite struct {
	app     *app.App
	baseURL string
	client  *resty.Client
}

var suite *RersTestSuite

var _ = BeforeSuite(func() {
	suite = new(RersTestSuite)
	err := suite.setup()
	Expect(err).NotTo(HaveOccurred())
})

// var _ = AfterSuite(func() {
// 	s.app.Context.Orm.Close()
// 	s.app.Srv.Shutdown()
// })

const (
	cfgFileName    = "./test-config/test-config.toml"
	fixturesFolder = "./fixtures"
)

// setup called once before test
func (s *RersTestSuite) setup() error {
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
func (s *RersTestSuite) setupServer() error {
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
	// wait til server started then return
	return s.waitServerStart(3 * time.Second)
}

//------------------------------------------------------------------------------
func (s *RersTestSuite) waitServerStart(timeout time.Duration) error {
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
