package twisk

import (
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/golang/protobuf/ptypes"
	"github.com/ribice/twisk/rpc/user"
)

// User represents user domain model
type User struct {
	ID                 int64
	FirstName          string
	LastName           string
	Username           string
	Password           string
	Email              string
	Phone              string
	Address            string
	Active             bool
	Token              string
	RoleID             int32
	TenantID           int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	LastLogin          *time.Time
	LastPasswordChange *time.Time
}

// AuthUser represents data stored in session/context for a user
type AuthUser struct {
	ID       int64
	TenantID int32
	Username string
	Email    string
	Role     AccessRole
}

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (u *User) BeforeInsert(_ orm.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (u *User) BeforeUpdate(_ orm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateLoginDetails updates login related fields
func (u *User) UpdateLoginDetails(token string) {
	u.Token = token
	t := time.Now()
	u.LastLogin = &t
}

// Delete sets deleted_at time to current_time
func (u *User) Delete() {
	t := time.Now()
	u.DeletedAt = &t
}

// Proto converts user db to proto User
func (u *User) Proto() *user.Resp {
	up := &user.Resp{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Address:   u.Address,
		Active:    u.Active,
		TenantId:  u.TenantID,
		RoleName:  user.Resp_RoleName(int32(u.RoleID)),
	}

	if !u.CreatedAt.IsZero() {
		createdAt, _ := ptypes.TimestampProto(u.CreatedAt)
		up.CreatedAt = createdAt
	}

	if !u.UpdatedAt.IsZero() {
		updatedAt, _ := ptypes.TimestampProto(u.UpdatedAt)
		up.UpdatedAt = updatedAt
	}

	if u.DeletedAt != nil {
		deletedAt, _ := ptypes.TimestampProto(*u.DeletedAt)
		up.DeletedAt = deletedAt
	}

	if u.LastLogin != nil {
		lastLogin, _ := ptypes.TimestampProto(*u.LastLogin)
		up.LastLogin = lastLogin
	}

	if u.LastPasswordChange != nil {
		lpc, _ := ptypes.TimestampProto(*u.LastPasswordChange)
		up.LastLogin = lpc
	}

	return up
}
