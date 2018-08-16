package mockdb

import (
	"github.com/go-pg/pg/orm"
	"github.com/ribice/twisk/model"
)

// User database mock
type User struct {
	CreateFn          func(orm.DB, twisk.User) (*twisk.User, error)
	ViewFn            func(orm.DB, int64) (*twisk.User, error)
	ListFn            func(orm.DB, string, int, int) ([]twisk.User, error)
	DeleteFn          func(orm.DB, *twisk.User) error
	UpdateFn          func(orm.DB, *twisk.User) (*twisk.User, error)
	FindByAuthFn      func(orm.DB, string) (*twisk.User, error)
	FindByTokenFn     func(orm.DB, string) (*twisk.User, error)
	UpdateLastLoginFn func(orm.DB, *twisk.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr twisk.User) (*twisk.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int64) (*twisk.User, error) {
	return u.ViewFn(db, id)
}

// List mock
func (u *User) List(db orm.DB, q string, limit, page int) ([]twisk.User, error) {
	return u.ListFn(db, q, limit, page)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr *twisk.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr *twisk.User) (*twisk.User, error) {
	return u.UpdateFn(db, usr)
}

// FindByAuth mock
func (u *User) FindByAuth(db orm.DB, auth string) (*twisk.User, error) {
	return u.FindByAuthFn(db, auth)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, auth string) (*twisk.User, error) {
	return u.FindByTokenFn(db, auth)
}

// UpdateLastLogin mock
func (u *User) UpdateLastLogin(db orm.DB, usr *twisk.User) error {
	return u.UpdateLastLoginFn(db, usr)
}
