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

// Generic responses for framework

type errorMessage struct {
	// Error message
	//
	// status: 400, 401, 403, 404, 405, 415. Details are in message.
	Message string
}

// error returning by framework
// swagger:response
type echoHTTPErrorResponse struct {
	// in: body
	Body errorMessage
}

// OK response (No content)
// swagger:response
type emptyOKResponse struct {
}
