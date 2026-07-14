package controllers

import (
	"testing"

	"teamitmivhs/work-order-backend/models"
)

func boolPtr(value bool) *bool { return &value }

func TestDeriveGuestPriority(t *testing.T) {
	tests := []struct {
		name string
		req  models.WorkOrderRequest
		want string
	}{
		{"forged priority is ignored", models.WorkOrderRequest{Priority: "high", ImpactScope: "individual", Disruption: "minor", WorkaroundAvailable: boolPtr(true)}, "low"},
		{"default medium", models.WorkOrderRequest{ImpactScope: "room", Disruption: "partial", WorkaroundAvailable: boolPtr(false)}, "medium"},
		{"high requires full stop", models.WorkOrderRequest{ImpactScope: "multiple", Disruption: "stopped", WorkaroundAvailable: boolPtr(false)}, "high"},
		{"single device stays medium", models.WorkOrderRequest{ImpactScope: "individual", Disruption: "stopped", WorkaroundAvailable: boolPtr(false)}, "medium"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := deriveGuestPriority(test.req)
			if err != nil || got != test.want {
				t.Fatalf("deriveGuestPriority() = %q, %v; want %q", got, err, test.want)
			}
		})
	}
}

func TestDeriveGuestPriorityRequiresCompleteTriage(t *testing.T) {
	if _, err := deriveGuestPriority(models.WorkOrderRequest{}); err == nil {
		t.Fatal("deriveGuestPriority() accepted incomplete triage")
	}
}
