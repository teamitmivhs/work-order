package repository

import (
	"errors"
	"testing"
)

func TestManualMemberStatusAllowed(t *testing.T) {
	for status, want := range map[string]bool{
		"standby":   true,
		"nextshift": false,
		"offduty":   true,
		"onjob":     false,
	} {
		if got := manualMemberStatusAllowed(status); got != want {
			t.Errorf("manualMemberStatusAllowed(%q) = %v, want %v", status, got, want)
		}
	}
}

func TestUpdateMemberStatusRejectsManualOnJob(t *testing.T) {
	repo := &memberRepository{}
	if err := repo.UpdateMemberStatus(1, "onjob"); !errors.Is(err, ErrWorkOrderManagedStatus) {
		t.Fatalf("UpdateMemberStatus(onjob) error = %v, want ErrWorkOrderManagedStatus", err)
	}
}
