package query

import (
	"fmt"

	"github.com/ribice/twisk/model"
)

// ForTenant returns query for filtering rows by tenant_id
func ForTenant(u *twisk.AuthUser, tenantID int32) string {
	switch u.Role {
	case twisk.SuperAdminRole, twisk.AdminRole:
		if tenantID != 0 {
			return fmt.Sprintf("tenant_id = %v", tenantID)
		}
		return ""
	default:
		return fmt.Sprintf("tenant_id = %v", u.TenantID)

	}
}

// Pagination holds paginations data
type Pagination struct {
	Limit  int32
	Offset int32
}

// Paginate returns pagination details with default query limit
func Paginate(limit, page int32) (int, int) {
	if limit < 1 {
		limit = 50
	}
	return int(limit), int(limit * (page - 1))

}
