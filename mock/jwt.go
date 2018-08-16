package mock

import (
	"github.com/ribice/twisk/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*twisk.AuthUser) (string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *twisk.AuthUser) (string, error) {
	return j.GenerateTokenFn(u)
}
