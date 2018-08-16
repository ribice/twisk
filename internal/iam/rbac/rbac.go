package rbac

import (
	"context"

	"github.com/ribice/twisk/model"
)

// Service is RBAC application service
type Service struct{}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c context.Context, r twisk.AccessRole) bool {
	role, ok := c.Value("role").(twisk.AccessRole)
	return ok && !(role > r)
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c context.Context, ID int64) bool {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	id, ok := c.Value("id").(int64)
	return ok && (id == ID || s.isAdmin(c))
}

// EnforceTenant checks whether the request to apply change to tenant data
// is done by the user belonging to the that tenant and that the user has role tenantAdmin.
// If user has admin role, the check for tenant doesn't need to pass.
func (s *Service) EnforceTenant(c context.Context, ID int32) bool {
	tenantID, ok := c.Value("tenant_id").(int32)
	return ok && (tenantID == ID || s.isAdmin(c))
}

func (s *Service) isAdmin(c context.Context) bool {
	role, ok := c.Value("role").(twisk.AccessRole)
	return ok && !(role > twisk.AdminRole)
}

// EnforceTenantAdmin checks tenant admin
func (s *Service) EnforceTenantAdmin(c context.Context, ID int32) bool {
	// Must query company ID in database for the given user
	tenantID, ok := c.Value("tenant_id").(int32)
	if !ok {
		return false
	}
	role, ok := c.Value("role").(twisk.AccessRole)
	return ok && ((!(role > twisk.TenantAdminRole) && tenantID == ID) || s.isAdmin(c))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c context.Context, r twisk.AccessRole) bool {
	role, ok := c.Value("role").(twisk.AccessRole)
	return ok && !(role >= r)
}

// EnforceTenantAndRole performs auth check for same tenant and lower role.
// Used for user creation, deletion etc.
func (s *Service) EnforceTenantAndRole(c context.Context, roleID twisk.AccessRole, tenantID int32) bool {
	tID, ok := c.Value("tenant_id").(int32)
	if !ok {
		return false
	}
	role, ok := c.Value("role").(twisk.AccessRole)
	return ok && !(role >= roleID) &&
		(tID == tenantID || s.isAdmin(c))
}
