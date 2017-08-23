// Package server ...
// Schemes: http
// Host: localhost:11011
// BasePath: /
// Version: 0.0.1
// License: MIT
// Contact: Marat Kagarmanov<mz3corvinus@gmail.com>
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
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

	// swagger:route GET /version version GetVersion
	//
	// responses:
	//   default: echoHTTPErrorResponse
	e.GET("/version", versionHandler.GetVersion)

	// swagger:route GET /groups groups GetGroups
	//
	// responses:
	//   default: echoHTTPErrorResponse
	e.GET("/groups", groupsHandler.GetGroups)

	// swagger:route POST /users users CreateUser
	//
	// responses:
	//   default: echoHTTPErrorResponse
	e.POST("/users", usersHandler.CreateUser)

	// swagger:route GET /users users GetUsers
	//
	// responses:
	//   default: echoHTTPErrorResponse
	e.GET("/users", usersHandler.GetUsers)

	// swagger:route PUT /users/{ID} users PutUser
	//
	// responses:
	//   default: echoHTTPErrorResponse
	e.PUT("/users/:id", usersHandler.PutUser)

	// swagger:route DELETE /users/{ID} users DeleteUser
	//
	// responses:
	//	 default: echoHTTPErrorResponse
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
