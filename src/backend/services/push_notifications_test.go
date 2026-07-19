package services

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFormatNewWorkOrderNotification(t *testing.T) {
	title, body := formatNewWorkOrderNotification("Laptop Guru", "Ruang Guru", "Tidak bisa login", "high")
	if title != "Work Order Baru • Prioritas Tinggi" {
		t.Fatalf("unexpected title: %q", title)
	}
	wantBody := "Perangkat: Laptop Guru\nLokasi: Ruang Guru\nKendala: Tidak bisa login"
	if body != wantBody {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestTruncatePushBody(t *testing.T) {
	got := truncatePushBody(strings.Repeat("kendala ", 80))
	if utf8.RuneCountInString(got) > maxPushBodyRunes {
		t.Fatalf("push body has %d runes", utf8.RuneCountInString(got))
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("truncated push body must end with ellipsis: %q", got)
	}
}
