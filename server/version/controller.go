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
// swagger:model versionModel
type ResponseBody struct {
	Result     string `json:"result"`
	Version    string `json:"version"`
	ServerTime int64  `json:"server_time"`
}

// GetVersion is a GET /version handler
// swagger:route GET /version version GetVersion
//
// Returns server version and time.
//
//	produces:
//		- application/json
//
//	schemes: http
//
//	responses:
//		200: versionResponse
func (h *Handler) GetVersion(c echo.Context) error {
	response := ResponseBody{
		Result:     "OK",
		Version:    h.C.Config.Version,
		ServerTime: time.Now().UTC().Unix(),
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}
