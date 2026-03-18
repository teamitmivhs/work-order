package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// FIX: timeTrackerBaseURL tidak lagi dibaca saat package init dengan:
//
//	var timeTrackerBaseURL = os.Getenv("TIME_TRACKER_URL")
//
// Sebelumnya jika env belum di-set saat startup, nilainya akan selalu "" selamanya.
// Sekarang dibaca saat fungsi dipanggil agar bisa di-set kapan saja.
func getTimeTrackerBaseURL() (string, error) {
	url := os.Getenv("TIME_TRACKER_URL")
	if url == "" {
		return "", fmt.Errorf("TIME_TRACKER_URL environment variable is not set")
	}
	return url, nil
}

// getInternalAPIKey membaca shared secret untuk autentikasi ke Rust engine.
// FIX: sebelumnya tidak ada autentikasi sama sekali ke Rust service.
func getInternalAPIKey() string {
	key := os.Getenv("INTERNAL_API_KEY")
	if key == "" {
		return "changeme-internal-key"
	}
	return key
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type StartTimerRequest struct {
	WorkOrderID uint64 `json:"work_order_id"`
	ExecutorID  uint64 `json:"executor_id"`
}

type StopTimerRequest struct {
	WorkOrderID uint64 `json:"work_order_id"`
}

// StopTimerResponse memetakan respons dari Rust engine saat timer dihentikan
// FIX: struct baru agar Go bisa membaca durasi yang dikembalikan Rust
// dan menyimpannya ke database sebagai WorkingHours
type StopTimerResponse struct {
	WorkOrderID     uint64 `json:"work_order_id"`
	StartedAt       int64  `json:"started_at"`
	StoppedAt       int64  `json:"stopped_at"`
	DurationSeconds int64  `json:"duration_seconds"`
}

// doPost adalah helper internal untuk mengirim POST request ke Rust engine
// dengan header autentikasi dan error handling yang konsisten
func doPost(url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// FIX: kirim internal API key agar Rust engine bisa memvalidasi caller
	req.Header.Set("X-Internal-Key", getInternalAPIKey())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// StartTimer memanggil Rust engine untuk mulai menghitung waktu pengerjaan.
// Dipanggil dari TakeOrder setelah order berhasil diambil.
// FIX: fungsi ini sekarang benar-benar dipanggil di workorder_repository.TakeOrder
func StartTimer(workOrderID uint64, executorID uint64) error {
	baseURL, err := getTimeTrackerBaseURL()
	if err != nil {
		return err
	}

	resp, err := doPost(baseURL+"/timer/start", StartTimerRequest{
		WorkOrderID: workOrderID,
		ExecutorID:  executorID,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("start timer failed with status: %d", resp.StatusCode)
	}

	return nil
}

// StopTimer memanggil Rust engine untuk menghentikan timer dan mendapatkan durasi.
// FIX: sekarang return (int64, error) — durasi dalam detik dikembalikan ke caller
// agar bisa disimpan sebagai WorkingHours di database.
// Dipanggil dari CompleteOrder setelah order selesai.
func StopTimer(workOrderID uint64) (int64, error) {
	baseURL, err := getTimeTrackerBaseURL()
	if err != nil {
		return 0, err
	}

	resp, err := doPost(baseURL+"/timer/stop", StopTimerRequest{
		WorkOrderID: workOrderID,
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("stop timer failed with status: %d", resp.StatusCode)
	}

	// FIX: decode response untuk mendapatkan durasi aktual
	var result StopTimerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode stop timer response: %w", err)
	}

	return result.DurationSeconds, nil
}
