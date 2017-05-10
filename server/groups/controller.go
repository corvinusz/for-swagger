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
type getGroupParams struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
}

// OK response
// swagger:response getGroupsResponse
type getGroupsResponse struct {
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
	params, err := h.getParams(c)
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
func (h *Handler) getParams(c echo.Context) (*getGroupParams, error) {
	var err error
	qs := c.QueryParams()
	params := new(getGroupParams)
	// get id
	params.ID, err = utils.GetUintParamFromURL(qs, "id", 0)
	if err != nil {
		return nil, err
	}

	params.Name = qs.Get("name")

	return params, nil
}
