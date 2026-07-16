package controllers

import "testing"

func TestNormalizeMemberName(t *testing.T) {
	for _, test := range []struct {
		input, want string
		valid       bool
	}{
		{" Alvaro ", "Alvaro", true},
		{"ab", "ab", false},
		{"", "", false},
	} {
		got, valid := normalizeMemberName(test.input)
		if got != test.want || valid != test.valid {
			t.Errorf("normalizeMemberName(%q) = %q, %v; want %q, %v", test.input, got, valid, test.want, test.valid)
		}
	}
}

func TestNormalizeStaffRole(t *testing.T) {
	for input, want := range map[string]string{
		"":         "Operator",
		"operator": "Operator",
		"admin":    "Admin",
		"guru":     "Guru",
		" Guru ":   "Guru",
		"unknown":  "",
	} {
		if got := normalizeStaffRole(input); got != want {
			t.Errorf("normalizeStaffRole(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeStaffBatchYear(t *testing.T) {
	for _, test := range []struct {
		role, batch, want string
		valid             bool
	}{
		{"Guru", "", "", true},
		{"Guru", "14", "", true},
		{"Operator", " 14 ", "14", true},
		{"Operator", "", "", false},
	} {
		got, valid := normalizeStaffBatchYear(test.role, test.batch)
		if got != test.want || valid != test.valid {
			t.Errorf("normalizeStaffBatchYear(%q, %q) = %q, %v; want %q, %v", test.role, test.batch, got, valid, test.want, test.valid)
		}
	}
}

func TestNormalizeStaffDivisionForRole(t *testing.T) {
	for _, test := range []struct {
		role, division, want string
		valid                bool
	}{
		{"Guru", "Programmer", "", true},
		{"Guru", "", "", true},
		{"Operator", "programmer", "Programmer", true},
		{"Operator", "unknown", "", false},
	} {
		got, valid := normalizeStaffDivisionForRole(test.role, test.division)
		if got != test.want || valid != test.valid {
			t.Errorf("normalizeStaffDivisionForRole(%q, %q) = %q, %v; want %q, %v", test.role, test.division, got, valid, test.want, test.valid)
		}
	}
}
