// Package users ...
// swagger:meta
package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/corvinusz/for-swagger/server/groups"
)

// Entity represents user
// swagger:model userEntity
type Entity struct {
	ID       uint64        `json:"id" xorm:"'id' pk index autoincr"`
	Created  uint64        `json:"created" xorm:"'created'"`
	Updated  uint64        `json:"updated" xorm:"'updated'"`
	Login    string        `json:"login" xorm:"index unique"`
	Email    string        `json:"email" xorm:"'email'"`
	Password string        `json:"-"`
	GroupID  uint64        `json:"group_id" xorm:"'group_id' index"`
	Group    groups.Entity `json:"group" xorm:"-"`
}

// TableName used by xorm to set table name for entity
func (e *Entity) TableName() string {
	return "users"
}

// FindByParams users in database
func FindByParams(orm *xorm.Engine, params *getUsersParams) ([]Entity, error) {
	var (
		users []Entity
		err   error
	)
	err = orm.Asc("id").
		Where("login LIKE ?", params.Login+"%").
		And("email LIKE ?", params.Email+"%").
		Find(&users, &Entity{ID: params.ID, GroupID: params.GroupID})

	if err != nil {
		return users, err
	}

	// extract groups
	for i := range users {
		orm.ID(users[i].GroupID).Get(&users[i].Group)
	}

	return users, err
}

// ExtractFrom extracts user from database
func (e *Entity) ExtractFrom(orm *xorm.Engine) (int, error) {
	var status int

	found, err := orm.Get(e)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("not found")
	}
	// extract  subobjects
	if status, err = e.extractGroup(orm); err != nil {
		return status, err
	}

	return http.StatusOK, nil
}

// Save user to database
func (e *Entity) Save(orm *xorm.Engine) (int, error) {
	var (
		err      error
		status   int
		hash     []byte
		affected int64
	)
	// check if always exists
	affected, err = orm.Where("login = ?", e.Login).Count(&Entity{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected != 0 {
		return http.StatusConflict, errors.New("such login always exists")
	}
	// handle user-group foreign key
	if status, err = e.extractGroup(orm); err != nil {
		return status, err
	}
	// encrypt password
	hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	e.Password = string(hash[:])
	// set CreatedAt
	e.Created = uint64(time.Now().UTC().Unix())
	e.Updated = e.Created
	// save to DB
	affected, err = orm.InsertOne(e)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected == 0 {
		return http.StatusServiceUnavailable, errors.New("db refused to insert user")
	}
	return http.StatusCreated, nil
}

// Update user in database
func (e *Entity) Update(orm *xorm.Engine) (int, error) {
	var (
		err      error
		affected int64
		hash     []byte
	)
	// get old user
	old := &Entity{ID: e.ID}
	status, err := old.ExtractFrom(orm)
	if err != nil {
		return status, err
	}
	// process user-group foreign key
	if e.GroupID == 0 {
		e.GroupID = old.GroupID
	}
	if status, err = e.extractGroup(orm); err != nil {
		return status, err
	}
	// encrypt password
	if len(e.Password) != 0 {
		hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusServiceUnavailable, err
		}
		e.Password = string(hash[:])
	}
	// set updatedAt
	e.Updated = uint64(time.Now().UTC().Unix())
	// update
	affected, err = orm.ID(e.ID).Update(e)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected == 0 {
		return http.StatusServiceUnavailable, errors.New("db refused to update user")
	}
	if e.Created == 0 {
		e.Created = old.Created
	}
	return http.StatusOK, nil
}

// Delete user from database
func (e *Entity) Delete(orm *xorm.Engine) (int, error) {
	var (
		err   error
		found bool
		user  Entity
	)
	// check if user exists
	found, err = orm.ID(e.ID).Get(&user)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("user not found")
	}
	//delete
	_, err = orm.ID(e.ID).Delete(&Entity{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil
}

//------------------------------------------------------------------------------
func (e *Entity) extractGroup(orm *xorm.Engine) (int, error) {
	if e.GroupID == 0 {
		return http.StatusOK, nil
	}
	// extract license
	e.Group.ID = e.GroupID
	return e.Group.ExtractFrom(orm)
}
