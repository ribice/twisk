package pgsql

import (
	"fmt"
	"strings"

	"github.com/ribice/twisk/model"

	"github.com/go-pg/pg/orm"
)

// NewUser returns a new User instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// Create creates new user
func (s *User) Create(cl orm.DB, req twisk.User) (*twisk.User, error) {
	var user = new(twisk.User)
	res, err := cl.Query(user, "select id from users where (lower(username) = ? or lower(email) = ?) and deleted_at is null",
		strings.ToLower(req.Username), strings.ToLower(req.Email))

	if err != nil {
		return nil, err
	}

	if res.RowsReturned() != 0 {
		return nil, fmt.Errorf("username or email already exists")
	}

	if err := cl.Insert(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

// View returns single user by ID/PrimaryKey
func (s *User) View(cl orm.DB, id int64) (*twisk.User, error) {
	var user = new(twisk.User)
	if err := cl.Model(user).Where("deleted_at is null").Select(); err != nil {
		return nil, err
	}
	return user, nil
}

// List returns list of all users retreivable for the current user, depending on role
func (s *User) List(cl orm.DB, query string, limit, offset int) ([]twisk.User, error) {
	var users []twisk.User
	q := cl.Model(&users).Limit(limit).Offset(offset).Where("deleted_at is null").Order("id desc")
	if query != "" {
		q.Where(query)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return users, nil
}

// Delete delets a user by ID
func (s *User) Delete(cl orm.DB, user *twisk.User) error {
	user.Delete()
	_, err := cl.Model(user).Column("deleted_at").WherePK().Update()
	return err
}

// Update updates a single user
func (s *User) Update(cl orm.DB, user *twisk.User) (*twisk.User, error) {
	_, err := cl.Model(user).WherePK().Update()
	return user, err
}
