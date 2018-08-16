package pgsql

import (
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

// FindByAuth finds user by either username or email
func (s *User) FindByAuth(cl orm.DB, auth string) (*twisk.User, error) {
	var user = new(twisk.User)

	lAuth := strings.ToLower(auth)
	if err := cl.Model(user).
		Where("deleted_at is null").
		Where("lower(username) = ? or lower(email) = ?",
			lAuth, lAuth).Select(); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByToken finds user by either username or email
func (s *User) FindByToken(cl orm.DB, token string) (*twisk.User, error) {
	var user = new(twisk.User)

	if err := cl.Model(user).Where("token = ?", token).
		Where("deleted_at is null").Select(); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateLastLogin updates user's last login details
func (s *User) UpdateLastLogin(cl orm.DB, user *twisk.User) error {
	_, err := cl.Model(user).Column("last_login", "token").
		WherePK().Update()
	return err
}
