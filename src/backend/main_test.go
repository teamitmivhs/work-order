package main

import "testing"

func TestIsAllowedOrigin(t *testing.T) {
	tests := map[string]bool{
		"http://localhost:4323":       true,
		"http://127.0.0.1:8080":       true,
		"http://192.168.100.47:4323":  true,
		"https://192.168.100.47:8443": true,
		"https://evil.example":        false,
		"http://8.8.8.8:4323":         false,
		"http://192.168.100.47:3000":  false,
		"not-an-origin":               false,
	}
	for origin, want := range tests {
		if got := isAllowedOrigin(origin); got != want {
			t.Errorf("isAllowedOrigin(%q) = %v, want %v", origin, got, want)
		}
	}
}
