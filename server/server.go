// Package server Example of RESTful-server with swagger documentation
//
// The purpose of this application is to provide an exampe of API documentation
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
// Schemes: http
// Host: localhost:11011
// BasePath: /
// Version: 0.0.1
// License: MIT
// Contact: Marat Kagarmanov<mz3corvinus@gmail.com>
// swagger:meta
package server

import (
	stdctx "context"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/corvinusz/for-swagger/ctx"
	"github.com/corvinusz/for-swagger/server/groups"
	"github.com/corvinusz/for-swagger/server/users"
	"github.com/corvinusz/for-swagger/server/version"
)

// Server ...
type Server struct {
	context *ctx.Context
	echo    *echo.Echo
}

// New is a constructor
func New(c *ctx.Context) *Server {
	s := new(Server)
	s.context = c
	return s
}

// Start registers web API and starts http-server
func (s *Server) Start() {
	// create 'Echo' instance
	e := echo.New()
	s.echo = e
	//e.Logger.SetLevel(log.ERROR) // internal echo logger still always active

	// Global Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	var (
		versionHandler = version.Handler{C: s.context}
		groupsHandler  = groups.Handler{C: s.context}
		usersHandler   = users.Handler{C: s.context}
	)

	// version routes
	e.GET("/", versionHandler.GetVersion)
	e.GET("/version", versionHandler.GetVersion)
	// groups
	e.GET("/groups", groupsHandler.GetGroups)
	// users
	e.POST("/users", usersHandler.CreateUser)
	e.GET("/users", usersHandler.GetUsers)
	e.PUT("/users/:id", usersHandler.PutUser)
	e.DELETE("/users/:id", usersHandler.DeleteUser)

	// start server
	e.Server.Addr = ":" + s.context.Config.Port
	e.Server.WriteTimeout = 5 * time.Second
	e.Server.ReadTimeout = 5 * time.Second
	s.context.Logger.Println("starting server at " + e.Server.Addr)
	err := e.Start(e.Server.Addr)
	if err != nil {
		s.context.Logger.Println("server error: ", err.Error())
	}
}

// Shutdown ...
func (s *Server) Shutdown() error {
	ctx, cancel := stdctx.WithTimeout(stdctx.Background(), 15*time.Second)
	defer cancel()
	return s.echo.Shutdown(ctx)
}