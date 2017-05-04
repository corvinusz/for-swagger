// Package version ...
// swagger:meta
package version

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/corvinusz/for-swagger/ctx"
)

// Handler is a application context carrier
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

// VersionResponseEnvelop is a wrapper
// swagger:response
type VersionResponseEnvelop struct {
	// Version Response
	// in: body
	Body *ResponseBody
}

// GetVersion is a GET /version handler
// swagger:route GET /version GetVersion
//
// Returns server version and time.
//
// responses:
//		200: VersionResponseEnvelop
func (h *Handler) GetVersion(c echo.Context) error {
	response := ResponseBody{
		Result:     "OK",
		Version:    h.C.Config.Version,
		ServerTime: time.Now().UTC().Unix(),
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}
