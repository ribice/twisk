package iam

import (
	"context"

	"github.com/ribice/twisk/model"

	"github.com/rs/xid"

	"github.com/twitchtv/twirp"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ribice/twisk/rpc/iam"
)

// New instantiates new IAM service
func New(db *pg.DB, tg TokenGenerator, udb UserDB, sec Securer) *Service {
	return &Service{db: db, tg: tg, udb: udb, sec: sec}
}

// Service represents IAM application service
type Service struct {
	db  *pg.DB
	tg  TokenGenerator
	udb UserDB
	sec Securer
}

// TokenGenerator generates new jwt token
type TokenGenerator interface {
	GenerateToken(*twisk.AuthUser) (string, error)
}

// UserDB represents user database interface
type UserDB interface {
	FindByAuth(orm.DB, string) (*twisk.User, error)
	FindByToken(orm.DB, string) (*twisk.User, error)
	UpdateLastLogin(orm.DB, *twisk.User) error
}

// Securer represents password securing service
type Securer interface {
	MatchesHash(string, string) bool
}

var (
	invalidUserPW = twirp.NewError(twirp.PermissionDenied, "invalid username or password")
	invalidToken  = twirp.NewError(twirp.PermissionDenied, "invalid token")
)

// Auth tries to authenticate user given username and password
func (s *Service) Auth(c context.Context, req *iam.AuthReq) (*iam.AuthResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dbCtx := s.db.WithContext(c)

	usr, err := s.udb.FindByAuth(dbCtx, req.Auth)
	if err != nil {
		return nil, invalidUserPW
	}

	if !s.sec.MatchesHash(usr.Password, req.Password) {
		return nil, invalidUserPW
	}

	token, err := s.tg.GenerateToken(&twisk.AuthUser{
		ID:       usr.ID,
		TenantID: usr.TenantID,
		Username: usr.Username,
		Email:    usr.Email,
		Role:     twisk.AccessRole(usr.RoleID),
	})

	if err != nil {
		return nil, err
	}

	uToken := xid.New().String()

	usr.UpdateLoginDetails(uToken)

	if err = s.udb.UpdateLastLogin(dbCtx, usr); err != nil {
		return nil, err
	}

	return &iam.AuthResp{
		Token:        token,
		RefreshToken: uToken,
	}, nil
}

// Refresh refreshes user's jwt token
func (s *Service) Refresh(c context.Context, req *iam.RefreshReq) (*iam.RefreshResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	usr, err := s.udb.FindByToken(s.db.WithContext(c), req.Token)
	if err != nil {
		return nil, invalidToken
	}

	token, err := s.tg.GenerateToken(&twisk.AuthUser{
		ID:       usr.ID,
		TenantID: usr.TenantID,
		Username: usr.Username,
		Email:    usr.Email,
		Role:     twisk.AccessRole(usr.RoleID),
	})

	if err != nil {
		return nil, err
	}

	return &iam.RefreshResp{
		Token: token,
	}, nil
}
