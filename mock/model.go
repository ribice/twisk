package mock

import (
	"errors"

	"github.com/go-pg/pg"
)

var (
	// ErrGeneric used for testing purposes
	ErrGeneric = errors.New("generic error used for testing purposes")
)

type bleja struct {
	p *pg.DB
}
