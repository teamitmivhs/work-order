package controllers

import "testing"

func TestNormalizeWeeklyShiftMembers(t *testing.T) {
	for _, test := range []struct {
		ids   []int
		valid bool
	}{
		{[]int{1, 2, 3}, true},
		{[]int{1, 1, 3}, false},
		{[]int{1, 2}, false},
		{[]int{1, 2, 0}, false},
	} {
		_, valid := normalizeWeeklyShiftMembers(test.ids)
		if valid != test.valid {
			t.Errorf("normalizeWeeklyShiftMembers(%v) valid = %v, want %v", test.ids, valid, test.valid)
		}
	}
}

func TestShiftDayName(t *testing.T) {
	if got := shiftDayName(5); got != "Jumat" {
		t.Fatalf("shiftDayName(5) = %q, want Jumat", got)
	}
	if got := shiftDayName(6); got != "" {
		t.Fatalf("shiftDayName(6) = %q, want empty", got)
	}
}
