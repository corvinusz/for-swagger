// Package server ...
// swagger:meta
package server

type errorMessage struct {
	// Error message
	//
	// http response status is one of [400, 401, 403, 404, 405, 415]
	Message string
}

// error returning by framework
//swagger:response
type echoHTTPErrorResponse struct {
	// in: body
	Body errorMessage
}

// OK response (No content)
//swagger:response
type echoOKResponse struct {
}
