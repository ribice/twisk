package mock

import (
	"context"

	"github.com/ribice/twisk/model"
)

// Auth mock
type Auth struct {
	GetUserFn func(context.Context) *twisk.AuthUser
}

// GetUser mock
func (s *Auth) GetUser(c context.Context) *twisk.AuthUser {
	return s.GetUserFn(c)
}
