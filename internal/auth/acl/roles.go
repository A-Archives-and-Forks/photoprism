package acl

import (
	"sort"
	"strings"
)

// RoleStrings represents user role names mapped to roles.
type RoleStrings map[string]Role

// UserRoles maps valid user account roles.
var UserRoles = RoleStrings{
	string(RoleAdmin):   RoleAdmin,
	string(RoleGuest):   RoleGuest,
	string(RoleVisitor): RoleVisitor,
	string(RoleNone):    RoleNone,
	RoleAliasNone:       RoleNone,
}

// ClientRoles maps valid API client roles.
var ClientRoles = RoleStrings{
	string(RoleAdmin):    RoleAdmin,
	string(RoleInstance): RoleInstance,
	"app":                RoleInstance,
	string(RoleService):  RoleService,
	string(RolePortal):   RolePortal,
	string(RoleClient):   RoleClient,
	string(RoleNone):     RoleNone,
	RoleAliasNone:        RoleNone,
}

// AdminRoles maps the roles that grant administrative privileges. The
// Portal-only cluster_admin is treated as an admin-tier role everywhere admin
// privileges are checked (e.g. user-management self-lockout protection), so a
// cluster_admin owner is not forced or downgraded to the plain admin role.
var AdminRoles = RoleStrings{
	string(RoleAdmin):        RoleAdmin,
	string(RoleClusterAdmin): RoleClusterAdmin,
}

// IsAdminRole reports whether role is an administrative role (admin or cluster_admin).
func IsAdminRole(role Role) bool {
	_, ok := AdminRoles[string(role)]
	return ok
}

// IsFederatedRole reports whether role may be assigned to a user account through
// an external identity provider (OIDC group/role claims or LDAP group/attribute
// mapping). The Portal-only operator role cluster_admin and the anonymous
// visitor role are never federated, and the empty role is rejected so a missing
// mapping cannot clear an account's role. Federation therefore neither grants
// nor revokes a non-federatable role: a compromised IdP/AD cannot escalate an
// account to operator access (the Portal Admin UI), and an existing
// cluster_admin/visitor account is never changed by a directory sync.
func IsFederatedRole(role Role) bool {
	switch role {
	case RoleNone, RoleClusterAdmin, RoleVisitor:
		return false
	default:
		return true
	}
}

// FederatedRoleUpdate reports the account role an external identity provider may
// apply to an existing user and whether to apply it. Federation neither grants
// nor revokes a non-federatable role, so it returns ok=false when the current
// role is non-federatable (cluster_admin/visitor must not be touched), when the
// mapped role is non-federatable (no escalation to operator/visitor), or when
// the role is unchanged. Both the OIDC and LDAP sync paths share this decision.
func FederatedRoleUpdate(current, mapped Role) (Role, bool) {
	if !IsFederatedRole(current) || !IsFederatedRole(mapped) || current == mapped {
		return RoleNone, false
	}

	return mapped, true
}

// Strings returns the roles as string slice for display, e.g. CLI help.
func (m RoleStrings) Strings() []string {
	result := make([]string, 0, len(m))

	for r := range m {
		if r == "" || r == RoleAliasNone || r == "app" || r == RoleVisitor.String() {
			continue
		}

		result = append(result, r)
	}

	sort.Strings(result)

	return result
}

// String returns the comma separated roles as string.
func (m RoleStrings) String() string {
	return strings.Join(m.Strings(), ", ")
}

// CliUsageString returns the roles as string for use in CLI usage descriptions.
func (m RoleStrings) CliUsageString() string {
	s := m.Strings()

	if l := len(s); l > 1 {
		s[l-1] = "or " + s[l-1]
	}

	return strings.Join(s, ", ")
}

// Roles grants permissions to roles.
type Roles map[Role]Grant

// Allow checks whether the permission is granted based on the role.
func (roles Roles) Allow(role Role, grant Permission) bool {
	if a, ok := roles[role]; ok {
		return a.Allow(grant)
	} else if a, ok = roles[RoleDefault]; ok {
		return a.Allow(grant)
	}

	return false
}
