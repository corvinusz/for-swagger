// Package version ...
// swagger:meta
package version

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/corvinusz/for-swagger/ctx"
)

// Handler is an application context carrier
type Handler struct {
	C *ctx.Context
}

// ResponseBody defines response body on GET /version
type ResponseBody struct {
	// Response Body
	Result     string `json:"result"`
	Version    string `json:"version"`
	ServerTime int64  `json:"server_time"`
}

// http response on GET /version
// swagger:response
type versionResponse struct {
	// response OK
	// in: body
	Body *ResponseBody
}

// GetVersion is a GET /version handler
// swagger:route GET /version version GetVersion
//
// Version Response
//
// responses:
//		200: versionResponse
//
//		default: echoHTTPErrorResponse
func (h *Handler) GetVersion(c echo.Context) error {
	response := ResponseBody{
		Result:     "OK",
		Version:    h.C.Config.Version,
		ServerTime: time.Now().UTC().Unix(),
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}
