package twisk

import "context"

// AccessRole represents access role type
type AccessRole int8

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = iota + 1

	// AdminRole has admin specific permissions
	AdminRole

	// TenantAdminRole can edit tenant specific things
	TenantAdminRole

	// UserRole is a standard user
	UserRole
)

// RBACService represents role-based access control service interface
type RBACService interface {
	EnforceRole(context.Context, AccessRole) bool
	EnforceUser(context.Context, int64) bool
	EnforceTenant(context.Context, int32) bool
	EnforceTenantAdmin(context.Context, int32) bool
	EnforceTenantAndRole(context.Context, AccessRole, int32) bool
	IsLowerRole(context.Context, AccessRole) bool
}

// Role entity
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
