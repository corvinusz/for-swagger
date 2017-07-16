package groups

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-xorm/xorm"
)

// Entity represents group
// swagger:model groupEntity
type Entity struct {
	ID      uint64 `json:"id" xorm:"'id' pk index autoincr"`
	Created uint64 `json:"created" xorm:"'created'"`
	Updated uint64 `json:"updated" xorm:"'updated'"`
	Name    string `json:"name" xorm:"index unique 'name'"`
	Note    string `json:"note"`
}

// TableName is used by xorm to set a table name for this entity
func (e *Entity) TableName() string {
	return "groups"
}

// FindByParams return groups that matches params
func FindByParams(orm *xorm.Engine, params *getGroupParams) ([]Entity, error) {
	var (
		groups []Entity
		err    error
	)
	err = orm.Asc("id").
		Where("name LIKE ?", params.Name+"%").
		Find(&groups, &Entity{ID: params.ID})

	return groups, err
}

// ExtractFrom extracts group from database with strict data
func (e *Entity) ExtractFrom(orm *xorm.Engine) (int, error) {
	found, err := orm.Get(e)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("group not found")
	}
	return http.StatusOK, nil
}

// Save group
func (e *Entity) Save(orm *xorm.Engine) (int, error) {
	var (
		err      error
		affected int64
	)
	// check if always exists
	affected, err = orm.Where("name = ?", e.Name).Count(&Entity{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected != 0 {
		return http.StatusConflict, errors.New("such name always exists")
	}
	// set CreatedAt
	e.Created = uint64(time.Now().UTC().Unix())
	e.Updated = e.Created
	// save to DB
	affected, err = orm.InsertOne(e)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected == 0 {
		return http.StatusServiceUnavailable, errors.New("db refused to insert group")
	}
	return http.StatusCreated, nil
}
