package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/ribice/twisk/model"

	"github.com/ribice/twisk/internal/user"
	"github.com/ribice/twisk/mock"
	"github.com/ribice/twisk/mock/mockdb"
	userpb "github.com/ribice/twisk/rpc/user"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		req      *userpb.CreateReq
		udb      *mockdb.User
		rbac     *mock.RBAC
		sec      *mock.Secure
		wantData *userpb.Resp
		wantErr  bool
	}{
		{
			name:    "Fail on validation",
			wantErr: true,
			req: &userpb.CreateReq{
				Username: "in__va__lid",
			},
		},
		{
			name:    "Fail on RBAC",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool {
					return false
				},
			},
			req: &userpb.CreateReq{
				Username:  "ribice",
				Email:     "ribice@gmail.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "testing",
				RoleId:    3,
				TenantId:  2,
			},
		},
		{
			name:    "Fail on insecure password",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool { return true },
			},
			sec: &mock.Secure{
				PasswordFn: func(string, ...string) bool { return false },
			},
			req: &userpb.CreateReq{
				Username:  "ribice",
				Email:     "ribice@gmail.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "testing",
				RoleId:    3,
				TenantId:  2,
			},
		},
		{
			name:    "Fail on creating user on db",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool { return true },
			},
			sec: &mock.Secure{
				PasswordFn: func(string, ...string) bool { return true },
				HashFn:     func(s string) string { return "ha$hed" },
			},
			udb: &mockdb.User{
				CreateFn: func(orm.DB, twisk.User) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.CreateReq{
				Username:  "ribice",
				Email:     "ribice@gmail.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "testing",
				RoleId:    3,
				TenantId:  2,
			},
		},
		{
			name: "Success",
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool { return true },
			},
			sec: &mock.Secure{
				PasswordFn: func(string, ...string) bool { return true },
				HashFn:     func(s string) string { return "ha$hed" },
			},
			udb: &mockdb.User{
				CreateFn: func(db orm.DB, u twisk.User) (*twisk.User, error) {
					return &u, nil
				},
			},
			req: &userpb.CreateReq{
				RoleId:    1,
				TenantId:  2,
				Email:     "large@mail.com",
				Username:  "juzernejm",
				Password:  "notHashed",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantData: &userpb.Resp{
				RoleName:  userpb.Resp_RoleName(int32(1)),
				TenantId:  2,
				Email:     "large@mail.com",
				Username:  "juzernejm",
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}
	dbcl := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(dbcl, tt.udb, tt.rbac, tt.sec, nil)
			usr, err := s.Create(nil, tt.req)
			if tt.wantData != nil {
				fmt.Printf("Case %v, User %v, Error %v", tt.name, usr, err)
				tt.wantData.CreatedAt = usr.CreatedAt
				tt.wantData.UpdatedAt = usr.UpdatedAt
			}
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		req      *userpb.ListReq
		wantData *userpb.ListResp
		udb      *mockdb.User
		auth     *mock.Auth
		wantErr  bool
	}{
		{
			name:    "Fail on validation",
			wantErr: true,
			req: &userpb.ListReq{
				Limit: 50, Page: -5,
			},
		},
		{
			name:    "Fail on UserDB list",
			wantErr: true,
			auth: &mock.Auth{
				GetUserFn: func(context.Context) *twisk.AuthUser {
					return &twisk.AuthUser{
						TenantID: 123,
						Role:     twisk.AdminRole,
					}
				},
			},
			udb: &mockdb.User{
				ListFn: func(orm.DB, string, int, int) ([]twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.ListReq{
				Limit: 50, Page: 2,
			},
		},
		{
			name: "Success",

			auth: &mock.Auth{
				GetUserFn: func(context.Context) *twisk.AuthUser {
					return &twisk.AuthUser{
						TenantID: 123,
						Role:     twisk.AdminRole,
					}
				},
			},
			udb: &mockdb.User{
				ListFn: func(orm.DB, string, int, int) ([]twisk.User, error) {
					return []twisk.User{
						{
							ID:        123,
							FirstName: "John",
							LastName:  "Doe",
							Email:     "johndoe@mail.com",
						},
						{
							ID:        321,
							FirstName: "Joanna",
							LastName:  "Doe",
							Email:     "joannadoe@mail.com",
						},
					}, nil
				},
			},
			req: &userpb.ListReq{
				Limit: 50, Page: 2,
			},
			wantData: &userpb.ListResp{
				Users: []*userpb.Resp{
					{
						ID:        123,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "johndoe@mail.com",
					},
					{
						ID:        321,
						FirstName: "Joanna",
						LastName:  "Doe",
						Email:     "joannadoe@mail.com",
					},
				}},
		},
	}
	dbcl := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(dbcl, tt.udb, nil, nil, tt.auth)
			usr, err := s.List(nil, tt.req)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name     string
		req      *userpb.IDReq
		wantData *userpb.Resp
		udb      *mockdb.User
		rbac     *mock.RBAC
		wantErr  bool
	}{
		{
			name:    "Fail on validation",
			wantErr: true,
			req: &userpb.IDReq{
				ID: -5,
			},
		},
		{
			name:    "Fail on userDB view",
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(orm.DB, int64) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
		},
		{
			name:    "Fail on rbac",
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
					}, nil
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
			rbac: &mock.RBAC{
				EnforceTenantFn: func(context.Context, int32) bool {
					return false
				},
			},
		},
		{
			name: "Success",
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
					}, nil
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
			rbac: &mock.RBAC{
				EnforceTenantFn: func(context.Context, int32) bool {
					return true
				},
			},
			wantData: &userpb.Resp{
				ID:        123,
				FirstName: "John",
				LastName:  "Doe",
				TenantId:  421,
			},
		},
	}
	dbcl := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(dbcl, tt.udb, tt.rbac, nil, nil)
			usr, err := s.View(nil, tt.req)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		req      *userpb.IDReq
		wantData *userpb.MessageResp
		udb      *mockdb.User
		rbac     *mock.RBAC
		wantErr  bool
	}{
		{
			name:    "Fail on validation",
			wantErr: true,
			req: &userpb.IDReq{
				ID: -5,
			},
		},
		{
			name:    "Fail on userDB view",
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(orm.DB, int64) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
		},
		{
			name:    "Fail on rbac",
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
						RoleID:    2,
					}, nil
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool {
					return false
				},
			},
		},
		{
			name:    "Fail on userDB delete",
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
						RoleID:    2,
					}, nil
				},
				DeleteFn: func(orm.DB, *twisk.User) error {
					return mock.ErrGeneric
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool {
					return true
				},
			},
		},
		{
			name: "Success",
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
						RoleID:    2,
					}, nil
				},
				DeleteFn: func(orm.DB, *twisk.User) error {
					return nil
				},
			},
			req: &userpb.IDReq{
				ID: 123,
			},
			rbac: &mock.RBAC{
				EnforceTenantAndRoleFn: func(context.Context, twisk.AccessRole, int32) bool {
					return true
				},
			},
			wantData: &userpb.MessageResp{
				Message: "OK",
			},
		},
	}
	dbcl := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(dbcl, tt.udb, tt.rbac, nil, nil)
			usr, err := s.Delete(nil, tt.req)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		req      *userpb.UpdateReq
		wantData *userpb.Resp
		udb      *mockdb.User
		rbac     *mock.RBAC
		wantErr  bool
	}{
		{
			name:    "Fail on validation",
			wantErr: true,
			req: &userpb.UpdateReq{
				FirstName: "Johnn",
				LastName:  "Doee",
			},
		},
		{
			name:    "Fail on Enforce User",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(context.Context, int64) bool {
					return false
				},
			},
			req: &userpb.UpdateReq{
				FirstName: "Johnn",
				LastName:  "Doee",
				ID:        123,
			},
		},
		{
			name:    "Fail on userDB view",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(context.Context, int64) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.UpdateReq{
				FirstName: "Johnn",
				LastName:  "Doee",
				ID:        123,
			},
		},
		{
			name:    "Fail on userDB update",
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(context.Context, int64) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
						RoleID:    2,
					}, nil
				},
				UpdateFn: func(orm.DB, *twisk.User) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			req: &userpb.UpdateReq{
				FirstName: "Johnn",
				LastName:  "Doee",
				ID:        123,
			},
		},
		{
			name: "Success",
			rbac: &mock.RBAC{
				EnforceUserFn: func(context.Context, int64) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int64) (*twisk.User, error) {
					return &twisk.User{
						ID:        id,
						FirstName: "John",
						LastName:  "Doe",
						TenantID:  421,
						RoleID:    2,
						Phone:     "3876192",
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *twisk.User) (*twisk.User, error) {
					return usr, nil
				},
			},
			req: &userpb.UpdateReq{
				ID:        321,
				FirstName: "Johnn",
				LastName:  "Doee",
			},
			wantData: &userpb.Resp{
				ID:        321,
				FirstName: "Johnn",
				LastName:  "Doee",
				TenantId:  421,
				RoleName:  userpb.Resp_RoleName(int32(2)),
			},
		},
	}
	dbcl := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(dbcl, tt.udb, tt.rbac, nil, nil)
			usr, err := s.Update(nil, tt.req)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
