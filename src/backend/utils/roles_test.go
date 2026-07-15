package utils

import "testing"

func TestIsAdminRole(t *testing.T) {
	for role, want := range map[string]bool{
		"Admin":    true,
		"Guru":     true,
		" guru ":   true,
		"Operator": false,
		"Guest":    false,
	} {
		if got := IsAdminRole(role); got != want {
			t.Errorf("IsAdminRole(%q) = %v, want %v", role, got, want)
		}
	}
}
