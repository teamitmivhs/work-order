package models

import "database/sql"

type WorkOrder struct {
	ID              int      `json:"id"`
	Priority        string   `json:"priority"`
	Time            string   `json:"time"`
	Requester       any      `json:"requester"`
	Location        string   `json:"location"`
	Device          string   `json:"device"`
	Problem         string   `json:"problem"`
	Executors       []int    `json:"executors"`
	WorkingHours    *int     `json:"workingHours,omitempty"`
	Status          string   `json:"status"`
	SafetyChecklist []string `json:"safetyChecklist"`
}

type Member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
}
type Summary struct {
	TotalWorkOrders      int `json:"totalWorkOrders"`
	PendingWorkOrders    int `json:"pendingWorkOrders"`
	InProgressWorkOrders int `json:"inProgressWorkOrders"`
	CompletedWorkOrders  int `json:"completedWorkOrders"`
}
type Kaizen struct {
	TotalKaizens       int `json:"totalKaizens"`
	ImplementedKaizens int `json:"implementedKaizens"`
	PendingKaizens     int `json:"pendingKaizens"`
}
type TechGuide struct {
	TotalArticles     int `json:"totalArticles"`
	PublishedArticles int `json:"publishedArticles"`
	DraftArticles     int `json:"draftArticles"`
}
type WorkOrderRequest struct {
	ID          int    `json:"id"`
	Priority    string `json:"priority"`
	TimeDisplay string `json:"time_display"`
	TimeSort    string `json:"time_sort"`
	Requester   string `json:"requester"`
	Location    string `json:"location"`
	Device      string `json:"device"`
	Problem     string `json:"problem"`

	// Executors harus diisi dari tabel relasi task_executors
	Executors []int `json:"executors"`

	// WorkingHours bisa NULL di DB
	WorkingHours sql.NullString `json:"workingHours"`
	Status       string         `json:"status"`

	// SafetyChecklist harus diisi dari tabel relasi task_safety_checklists
	SafetyChecklist []string `json:"safetyChecklist"`

	// CompletedAt bisa NULL di DB
	CompletedAt sql.NullString `json:"completedAt"`
}
type TakeWorkOrder struct {
	Executors            []int64  `json:"executors"`
	SafetyChecklistItems []string `json:"safety_checklist_items"`
	Status               string   `json:"status"`
}
type CompleteWorkOrder struct {
	Status             string `json:"status"`
	CompletedAtDisplay string `json:"completed_at_display"`
}
