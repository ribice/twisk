package user

import (
	"context"
	"strings"

	"github.com/ribice/twisk/internal/pkg/query"
	"github.com/ribice/twisk/internal/pkg/structs"
	"github.com/ribice/twisk/model"

	"github.com/twitchtv/twirp"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ribice/twisk/rpc/user"
)

// New creates new user application service
func New(db *pg.DB, udb DB, rbac twisk.RBACService, sec Securer, auth AuthService) *Service {
	return &Service{dbcl: db, udb: udb, rbac: rbac, sec: sec, auth: auth}
}

// Service represents user application service
type Service struct {
	dbcl *pg.DB
	udb  DB
	rbac twisk.RBACService
	auth AuthService
	sec  Securer
}

var (
	unauthorizedErr = twirp.NewError(twirp.PermissionDenied, "unathorized")
)

// DB represents user database interface
type DB interface {
	Create(orm.DB, twisk.User) (*twisk.User, error)
	View(orm.DB, int64) (*twisk.User, error)
	List(orm.DB, string, int, int) ([]twisk.User, error)
	Delete(orm.DB, *twisk.User) error
	Update(orm.DB, *twisk.User) (*twisk.User, error)
}

// AuthService represents authentication context service
type AuthService interface {
	GetUser(context.Context) *twisk.AuthUser
}

// Securer represents password securing service
type Securer interface {
	Hash(string) string
	Password(string, ...string) bool
}

// Create creates a new user account
func (s *Service) Create(c context.Context, req *user.CreateReq) (*user.Resp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !s.rbac.EnforceTenantAndRole(c, twisk.AccessRole(req.RoleId), req.TenantId) {
		return nil, unauthorizedErr
	}

	if !s.sec.Password(req.Password, req.FirstName, req.LastName, req.Email) {
		return nil, twirp.InternalError("password is not secure enough")
	}

	usr, err := s.udb.Create(s.dbcl.WithContext(c), twisk.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     strings.ToLower(req.Email),
		Password:  s.sec.Hash(req.Password),
		TenantID:  req.TenantId,
		RoleID:    req.RoleId,
	})

	if err != nil {
		return nil, err
	}

	return usr.Proto(), nil
}

// List returns list of users
func (s *Service) List(c context.Context, req *user.ListReq) (*user.ListResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	u := s.auth.GetUser(c)

	limit, offset := query.Paginate(req.Limit, req.Page)

	users, err := s.udb.List(
		s.dbcl.WithContext(c),
		query.ForTenant(u, req.TenantId),
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	var pu []*user.Resp
	for _, v := range users {
		pu = append(pu, v.Proto())
	}

	return &user.ListResp{Users: pu}, nil
}

// View returns single user
func (s *Service) View(c context.Context, req *user.IDReq) (*user.Resp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	usr, err := s.udb.View(s.dbcl.WithContext(c), req.ID)
	if err != nil {
		return nil, unauthorizedErr

	}

	if !s.rbac.EnforceTenant(c, usr.TenantID) {
		return nil, unauthorizedErr

	}

	return usr.Proto(), nil
}

// Delete deletes a user
func (s *Service) Delete(c context.Context, req *user.IDReq) (*user.MessageResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dbCtx := s.dbcl.WithContext(c)
	usr, err := s.udb.View(dbCtx, req.ID)
	if err != nil {
		return nil, unauthorizedErr
	}

	if !s.rbac.EnforceTenantAndRole(c, twisk.AccessRole(usr.RoleID), usr.TenantID) {
		return nil, unauthorizedErr
	}

	if err := s.udb.Delete(dbCtx, usr); err != nil {
		return nil, err
	}

	return &user.MessageResp{
		Message: "OK",
	}, nil

}

// Update updates user's contact information
func (s *Service) Update(c context.Context, req *user.UpdateReq) (*user.Resp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dbCtx := s.dbcl.WithContext(c)

	if !s.rbac.EnforceUser(c, req.ID) {
		return nil, unauthorizedErr
	}

	usr, err := s.udb.View(dbCtx, req.ID)
	if err != nil {
		return nil, unauthorizedErr
	}

	structs.Merge(usr, req)

	userUpdate, err := s.udb.Update(dbCtx, usr)
	if err != nil {
		return nil, err
	}

	return userUpdate.Proto(), nil
}
