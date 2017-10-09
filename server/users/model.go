package users

import (
	"net/http"
	"time"

	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/corvinusz/for-swagger/errors"
	"github.com/corvinusz/for-swagger/server/groups"
)

// Entity represents user
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
func (e *Entity) ExtractFrom(orm *xorm.Engine) error {
	found, err := orm.Get(e)
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	if !found {
		return errors.WrapCode(err, http.StatusNotFound)
	}
	// extract  subobjects
	err = e.extractGroup(orm)
	if err != nil {
		return err
	}
	return nil
}

// Save user to database
func (e *Entity) Save(orm *xorm.Engine) error {
	var (
		err      error
		hash     []byte
		affected int64
	)
	// check if always exists
	affected, err = orm.Where("login = ?", e.Login).Count(&Entity{})
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	if affected != 0 {
		return errors.NewWithCode("such login always exists", http.StatusConflict)
	}
	// handle user-group foreign key
	if err = e.extractGroup(orm); err != nil {
		return err
	}
	// encrypt password
	hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	e.Password = string(hash[:])
	// set CreatedAt
	e.Created = uint64(time.Now().UTC().Unix())
	e.Updated = e.Created
	// save to DB
	affected, err = orm.InsertOne(e)
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	if affected == 0 {
		return errors.NewWithCode("db refused to save user", http.StatusServiceUnavailable)
	}
	return nil
}

// Update user in database
func (e *Entity) Update(orm *xorm.Engine) error {
	var (
		affected int64
		hash     []byte
	)
	// get old user
	old := &Entity{ID: e.ID}
	err := old.ExtractFrom(orm)
	if err != nil {
		return err
	}
	// process user-group foreign key
	if e.GroupID == 0 {
		e.GroupID = old.GroupID
	}
	err = e.extractGroup(orm)
	if err != nil {
		return err
	}
	// encrypt password
	if len(e.Password) != 0 {
		hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.WrapCode(err, http.StatusServiceUnavailable)
		}
		e.Password = string(hash[:])
	}
	// set updatedAt
	e.Updated = uint64(time.Now().UTC().Unix())
	// update
	affected, err = orm.ID(e.ID).Update(e)
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	if affected == 0 {
		return errors.NewWithCode("db refused to update user", http.StatusServiceUnavailable)
	}
	if e.Created == 0 {
		e.Created = old.Created
	}
	return nil
}

// Delete user from database
func (e *Entity) Delete(orm *xorm.Engine) error {
	var (
		err   error
		found bool
		user  Entity
	)
	// check if user exists
	found, err = orm.ID(e.ID).Get(&user)
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	if !found {
		return errors.NewWithCode("user not found", http.StatusNotFound)
	}
	//delete
	_, err = orm.ID(e.ID).Delete(&Entity{})
	if err != nil {
		return errors.WrapCode(err, http.StatusServiceUnavailable)
	}
	return nil
}

//------------------------------------------------------------------------------
func (e *Entity) extractGroup(orm *xorm.Engine) error {
	if e.GroupID == 0 {
		return nil
	}
	// extract license
	e.Group.ID = e.GroupID
	return e.Group.ExtractFrom(orm)
}
