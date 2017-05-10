// Package groups ...
// swagger:meta
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

// swagger:parameters GetGroups
//in:query
type reqParams struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
}

// GET /groups response
// swagger:response getGroupsResponse
type getGroupsResponse struct {
	// response OK
	// in: body
	Groups []Entity
}

// GetGroups is a GET /groups handler
// swagger:operation GET /groups groups GetGroups
//
// Returns groups list
//
// ---
// responses:
//   '200':
//     "$ref": "#/responses/getGroupsResponse"
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
