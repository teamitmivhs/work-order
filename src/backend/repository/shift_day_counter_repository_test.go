package repository

import (
	"testing"
	"time"
)

func TestCrossedSunday(t *testing.T) {
	date := func(value string) time.Time {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			t.Fatal(err)
		}
		return parsed
	}
	for _, test := range []struct {
		last, current string
		want          bool
	}{
		{"2026-07-24", "2026-07-25", false},
		{"2026-07-25", "2026-07-26", true},
		{"2026-07-25", "2026-07-27", true},
		{"2026-07-26", "2026-07-27", false},
		{"2026-07-20", "2026-07-26", true},
		{"2026-07-26", "2026-07-26", false},
	} {
		if got := crossedSunday(date(test.last), date(test.current)); got != test.want {
			t.Errorf("crossedSunday(%s, %s) = %v, want %v", test.last, test.current, got, test.want)
		}
	}
}
