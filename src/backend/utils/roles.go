package utils

import "strings"

// IsAdminRole reports whether a role has full administrative access.
func IsAdminRole(role string) bool {
	role = strings.TrimSpace(role)
	return strings.EqualFold(role, "Admin") || strings.EqualFold(role, "Guru")
}
