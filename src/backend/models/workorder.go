package models

type WorkOrder struct {
	ID              int      `json:"id"`
	Priority        string   `json:"priority"`
	Time            string   `json:"time"`
	Requester       any      `json:"requester"`
	Location        string   `json:"location"`
	Device          string   `json:"device"`
	Problem         string   `json:"problem"`
	Executors       []int    `json:"executors"`
	WorkingHours    *int   `json:"workingHours,omitempty"`
	Status          string   `json:"status"`
	SafetyChecklist []string `json:"safetyChecklist"`
}

type Member struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Status string `json:"status"`
	Avatar string `json:"avatar"`
}
