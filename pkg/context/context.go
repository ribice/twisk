package context

import (
	"context"

	"github.com/ribice/twisk/model"
)

// KeyString should be used when setting and fetching context values
type KeyString string

// JWTKey is a context key for storing token
var JWTKey = "http_jwt_key"

// Service represents context service
type Service struct{}

// GetUser fetches auth user from context
func (s *Service) GetUser(c context.Context) *twisk.AuthUser {
	return c.Value(KeyString("_authuser")).(*twisk.AuthUser)
}
