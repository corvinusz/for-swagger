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

func TestApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bddtests Suite")
}

// Sui is testing context for suite
type Sui struct {
	app     *app.App
	baseURL string
	client  *resty.Client
}

var suite *Sui

var _ = BeforeSuite(func() {
	suite = new(Sui)
	err := suite.setup()
	Expect(err).NotTo(HaveOccurred())
})

const (
	cfgFileName    = "../test-config/test-config.toml"
	fixturesFolder = "../fixtures/initial"
)

// setup called once before test
func (s *Sui) setup() error {
	err := s.setupServer()
	if err != nil {
		return err
	}
	s.baseURL = "http://localhost:" + s.app.C.Config.Port
	// create and setup resty client
	s.client = resty.DefaultClient
	s.client.SetHeader("Content-Type", "application/json")
	s.client.SetHostURL(s.baseURL)
	return nil
}

//------------------------------------------------------------------------------
func (s *Sui) setupServer() error {
	var err error
	// init test application
	s.app, err = app.New(app.CmdFlags{CfgFileName: cfgFileName})
	if err != nil {
		return err
	}
	// load fixtures
	err = testfixtures.LoadFixtures(fixturesFolder, s.app.C.Orm.DB().DB, &testfixtures.SQLite{})
	if err != nil {
		return err
	}
	// start test server with go routine
	go s.app.Run()
	// wait til server started then return
	return s.waitServerStart(3 * time.Second)
}

//------------------------------------------------------------------------------
func (s *Sui) waitServerStart(timeout time.Duration) error {
	done := time.Now().Add(timeout)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", ":"+s.app.C.Config.Port)
		if err == nil {
			return c.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("cannot connect %v for %v", s.baseURL, timeout)
}
