// Package version ...
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

// GetVersionBody on GET /version
type GetVersionBody struct {
	Result     string `json:"result"`
	Version    string `json:"version"`
	ServerTime int64  `json:"server_time"`
}

// GetVersion is a GET /version handler
// swagger:operation GET /version version GetVersion
//
// Returns server time and version, also can be used as a healthcheck
//
// ---
// responses:
//   200:
//     description: Get server version.
//     schema:
//       $ref: '#/definitions/GetVersionBody'
func (h *Handler) GetVersion(c echo.Context) error {
	response := GetVersionBody{
		Result:     "OK",
		Version:    h.C.Config.Version,
		ServerTime: time.Now().UTC().Unix(),
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}
