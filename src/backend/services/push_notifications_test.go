package services

import "testing"

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
