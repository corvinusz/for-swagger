// Package users ...
// swagger:meta
package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/corvinusz/for-swagger/ctx"
	"github.com/corvinusz/for-swagger/utils"
)

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// Input represents payload data format
type Input struct {
	Login         string `json:"login"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	PasswordEtime uint64 `json:"password_etime"`
	GroupID       uint64 `json:"group_id"`
}

//------------------------------------------------------------------------------
// really not a json: json tags used for documenting query params
// swagger:parameters GetUsers
// in:query
type getUsersParams struct {
	Limit   uint64 `json:"limit"`
	Offset  uint64 `json:"offset"`
	ID      uint64 `json:"id"`
	GroupID uint64 `json:"group_id"`
	Login   string `json:"login"`
	Email   string `json:"email"`
}

// OK response
// swagger:response getUsersResponse
type getUsersResponse struct {
	// in: body
	Users []Entity
}

// GetUsers is a GET /users handler
// swagger:operation GET /users users GetUsers
//
// Returns users list
//
// ---
// responses:
//   '200':
//     "$ref": "#/responses/getUsersResponse"
func (h *Handler) GetUsers(c echo.Context) error {
	// parse parameters
	params, err := h.getQueryParams(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// extract users
	users, err := FindByParams(h.C.Orm, params)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

//------------------------------------------------------------------------------
// swagger:parameters CreateUser
type postUsersParams struct {
	//in:body
	Body *Input
}

// CREATED response
// swagger:response postUsersResponse
type postUsersResponse struct {
	// in: body
	User Entity
}

// CreateUser is a POST /users handler
// swagger:operation POST /users users CreateUser
//
// Create new user and save it in storage
//
// ---
// responses:
//   '201':
//     "$ref": "#/responses/postUsersResponse"
func (h *Handler) CreateUser(c echo.Context) error {
	var (
		status int
		err    error
		user   Entity
		input  Input
	)

	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate input
	if len(input.Login) == 0 {
		return c.String(http.StatusBadRequest, "bad login")
	}
	if len(input.Password) == 0 {
		return c.String(http.StatusBadRequest, "bad password")
	}
	if input.GroupID == 0 {
		return c.String(http.StatusBadRequest, "bad group_id")
	}
	// save
	user = Entity{
		Login:    input.Login,
		Password: input.Password,
		Email:    input.Email,
		GroupID:  input.GroupID,
	}

	status, err = user.Save(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(status, user)
}

//------------------------------------------------------------------------------
// swagger:parameters PutUser
type putUsersParams struct {
	//in:path
	ID uint64
	//in:body
	Body *Input
}

// OK response
// swagger:response putUsersResponse
type putUsersResponse struct {
	// in: body
	User Entity
}

// PutUser is a PUT /users/{ID} handler
// swagger:operation PUT /users/{ID} users PutUser
//
// Update user by id
//
// ---
// responses:
//   '200':
//     "$ref": "#/responses/putUsersResponse"
func (h *Handler) PutUser(c echo.Context) error {
	var (
		input  Input
		user   Entity
		id     uint64
		err    error
		status int
	)
	// parse id
	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// parse request body
	if err = c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// construct user
	user = Entity{
		ID:       id,
		Login:    input.Login,
		Email:    input.Email,
		GroupID:  input.GroupID,
		Password: input.Password,
	}
	// update
	status, err = user.Update(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

//------------------------------------------------------------------------------
// swagger:parameters DeleteUser
type deleteUsersParams struct {
	//in:path
	ID uint64
}

// DeleteUser is a DELETE /users/{ID} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var (
		id     uint64
		status int
		err    error
		user   Entity
	)

	id, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user.ID = id
	// delete
	status, err = user.Delete(h.C.Orm)
	if err != nil {
		return c.String(status, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

//------------------------------------------------------------------------------
func (h *Handler) getQueryParams(c echo.Context) (*getUsersParams, error) {
	var err error
	qs := c.QueryParams()
	params := new(getUsersParams)
	// get id
	params.ID, err = utils.GetUintParamFromURL(qs, "id", 0)
	if err != nil {
		return nil, err
	}
	// get group_id
	params.GroupID, err = utils.GetUintParamFromURL(qs, "group_id", 0)
	if err != nil {
		return nil, err
	}

	params.Login = c.QueryParam("login")
	params.Email = c.QueryParam("email")
	return params, nil
}
