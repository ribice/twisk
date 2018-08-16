package hooks

import (
	"context"
	"strings"

	"github.com/ribice/twisk/model"

	pkgctx "github.com/ribice/twisk/pkg/context"
	"github.com/twitchtv/twirp"
)

// TokenParser represents jwt auth token parser and validator
type TokenParser interface {
	ParseToken(string) (*twisk.AuthUser, error)
}

// WithJWTAuth creates new twirp authentication hook
func WithJWTAuth(parser TokenParser, skip ...string) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			mtd, _ := twirp.MethodName(ctx)
			for _, smtd := range skip {
				if mtd == smtd {
					return ctx, nil
				}
			}

			bearer, ok := ctx.Value(pkgctx.KeyString("HTTP-Authorization")).(string)
			if !ok {
				return ctx, twirp.NewError(twirp.Unauthenticated, "no auth headers present")
			}

			slice := strings.Split(bearer, " ")

			if len(slice) != 2 || strings.ToLower(slice[0]) != "bearer" {
				return ctx, twirp.NewError(twirp.Unauthenticated, "no bearer token")
			}

			user, err := parser.ParseToken(slice[1])
			if err != nil {
				return ctx, twirp.NewError(twirp.Unauthenticated, err.Error())
			}

			ctx = context.WithValue(ctx, pkgctx.KeyString(pkgctx.JWTKey), slice[1])
			return context.WithValue(ctx, pkgctx.KeyString("_authuser"), user), nil
		},
	}
}
