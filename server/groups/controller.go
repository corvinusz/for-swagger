package groups

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/corvinusz/for-swagger/ctx"
	"github.com/corvinusz/for-swagger/utils"
)

// Handler is a application context carrier
type Handler struct {
	C *ctx.Context
}

// Response defines http response on GET /groups
type Response struct {
	Groups  []Entity
	Remains uint64
}

// http request query string parameters
type reqParams struct {
	Limit  uint64
	Offset uint64
	ID     uint64
	Name   string
}

// GetGroups is a GET /groups handler
/* overcommented // swagger:route GET /groups groups GetGroups
//
// Returns user groups
//
// 	produces:
// 		- application/json
//
// 	schemes: http
//
// 	responses:
// 		200: Response*/
func (h *Handler) GetGroups(c echo.Context) error {
	// parse parameters
	params, err := h.getReqParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// extract groups
	groups, err := FindByParams(h.C.Orm, params)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	return c.JSON(http.StatusOK, groups)
}

//------------------------------------------------------------------------------
func (h *Handler) getReqParams(c echo.Context) (*reqParams, error) {
	var err error
	qs := c.QueryParams()
	params := new(reqParams)
	// get id
	params.ID, err = utils.GetUintParamFromURL(qs, "id", 0)
	if err != nil {
		return nil, err
	}

	params.Name = qs.Get("name")

	return params, nil
}
