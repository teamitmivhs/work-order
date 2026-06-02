package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var timeTrackerBaseURL = os.Getenv("TIME_TRACKER_URL")

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

func ensureBaseURL() error {
	if timeTrackerBaseURL == "" {
		return errors.New("TIME_TRACKER_URL is not set")
	}
	return nil
}

func StartTimer(workOrderID uint64, executorID uint64) error {
	if err := ensureBaseURL(); err != nil {
		return err
	}

	payload, err := json.Marshal(StartTimerRequest{
		WorkOrderID: workOrderID,
		ExecutorID:  executorID,
	})
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(
		timeTrackerBaseURL+"/timer/start",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("start timer failed, status: %d", resp.StatusCode)
	}

	return nil
}

func StopTimer(workOrderID uint64) error {
	if err := ensureBaseURL(); err != nil {
		return err
	}

	payload, err := json.Marshal(StopTimerRequest{
		WorkOrderID: workOrderID,
	})
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(
		timeTrackerBaseURL+"/timer/stop",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stop timer failed, status: %d", resp.StatusCode)
	}

	return nil
}
