package iam_test

import (
	"testing"

	"github.com/ribice/twisk/internal/iam"
	"github.com/ribice/twisk/model"

	"github.com/go-pg/pg/orm"

	"github.com/ribice/twisk/mock"

	"github.com/ribice/twisk/mock/mockdb"
	"github.com/stretchr/testify/assert"

	"github.com/go-pg/pg"
	iampb "github.com/ribice/twisk/rpc/iam"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		name     string
		req      *iampb.AuthReq
		udb      *mockdb.User
		sec      *mock.Secure
		tg       *mock.JWT
		wantErr  bool
		wantData *iampb.AuthResp
	}{
		{
			name: "Fail on validation",
			req: &iampb.AuthReq{
				Auth: "onlyauth",
			},
			wantErr: true,
		},
		{
			name: "Fail on FindByAuth",
			req: &iampb.AuthReq{
				Auth:     "email@mail.com",
				Password: "hunter2",
			},
			udb: &mockdb.User{
				FindByAuthFn: func(orm.DB, string) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on MatchesHash",
			req: &iampb.AuthReq{
				Auth:     "juzernejm",
				Password: "hunter2",
			},
			udb: &mockdb.User{
				FindByAuthFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						FirstName: "John",
						Password:  "(has*_*h3d)",
					}, nil
				},
			},
			sec: &mock.Secure{
				MatchesHashFn: func(string, string) bool {
					return false
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on GenerateToken",
			req: &iampb.AuthReq{
				Auth:     "juzernejm",
				Password: "hunter2",
			},
			udb: &mockdb.User{
				FindByAuthFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						FirstName: "John",
						Password:  "(has*_*h3d)",
					}, nil
				},
			},
			sec: &mock.Secure{
				MatchesHashFn: func(string, string) bool {
					return true
				},
			},
			tg: &mock.JWT{
				GenerateTokenFn: func(*twisk.AuthUser) (string, error) {
					return "", mock.ErrGeneric
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on UpdateLastLogin",
			req: &iampb.AuthReq{
				Auth:     "juzernejm",
				Password: "hunter2",
			},
			udb: &mockdb.User{
				FindByAuthFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						FirstName: "John",
						Password:  "(has*_*h3d)",
					}, nil
				},
				UpdateLastLoginFn: func(orm.DB, *twisk.User) error {
					return mock.ErrGeneric
				},
			},
			sec: &mock.Secure{
				MatchesHashFn: func(string, string) bool {
					return true
				},
			},
			tg: &mock.JWT{
				GenerateTokenFn: func(*twisk.AuthUser) (string, error) {
					return "jwttoken", nil
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			req: &iampb.AuthReq{
				Auth:     "juzernejm",
				Password: "hunter2",
			},
			udb: &mockdb.User{
				FindByAuthFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						FirstName: "John",
						Password:  "(has*_*h3d)",
					}, nil
				},
				UpdateLastLoginFn: func(orm.DB, *twisk.User) error {
					return nil
				},
			},
			sec: &mock.Secure{
				MatchesHashFn: func(string, string) bool {
					return true
				},
			},
			tg: &mock.JWT{
				GenerateTokenFn: func(*twisk.AuthUser) (string, error) {
					return "jwttoken", nil
				},
			},
			wantData: &iampb.AuthResp{
				Token: "jwttoken",
			},
		},
	}
	db := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := iam.New(db, tt.tg, tt.udb, tt.sec)
			resp, err := s.Auth(nil, tt.req)
			if tt.wantData != nil {
				tt.wantData.RefreshToken = resp.RefreshToken
			}
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestRefresh(t *testing.T) {
	cases := []struct {
		name     string
		req      *iampb.RefreshReq
		udb      *mockdb.User
		tg       *mock.JWT
		wantErr  bool
		wantData *iampb.RefreshResp
	}{
		{
			name: "Fail on validation",
			req: &iampb.RefreshReq{
				Token: "tooshort",
			},
			wantErr: true,
		},
		{
			name: "Fail on FindByToken",
			req: &iampb.RefreshReq{
				Token: "lengthis10lengthis20",
			},
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (*twisk.User, error) {
					return nil, mock.ErrGeneric
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on GenerateToken",
			req: &iampb.RefreshReq{
				Token: "lengthis10lengthis20",
			},
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						ID:       123,
						TenantID: 321,
						Username: "johndoe",
						Email:    "johndoe@mail.com",
						RoleID:   221,
					}, nil
				},
			},
			tg: &mock.JWT{
				GenerateTokenFn: func(*twisk.AuthUser) (string, error) {
					return "", mock.ErrGeneric
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			req: &iampb.RefreshReq{
				Token: "lengthis10lengthis20",
			},
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (*twisk.User, error) {
					return &twisk.User{
						ID:       123,
						TenantID: 321,
						Username: "johndoe",
						Email:    "johndoe@mail.com",
						RoleID:   221,
					}, nil
				},
			},
			tg: &mock.JWT{
				GenerateTokenFn: func(*twisk.AuthUser) (string, error) {
					return "newjwttoken", nil
				},
			},
			wantData: &iampb.RefreshResp{
				Token: "newjwttoken",
			},
		},
	}
	db := &pg.DB{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := iam.New(db, tt.tg, tt.udb, nil)
			resp, err := s.Refresh(nil, tt.req)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
