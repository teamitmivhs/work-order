package models

// WorkOrder adalah representasi response ke frontend
type WorkOrder struct {
	ID              int      `json:"id"`
	Priority        string   `json:"priority"`
	Time            string   `json:"time"`
	Requester       string   `json:"requester"` // FIX: any → string, selalu string dari DB
	Location        string   `json:"location"`
	Device          string   `json:"device"`
	Problem         string   `json:"problem"`
	Executors       []int    `json:"executors"`
	WorkingHours    *int     `json:"workingHours,omitempty"`
	Status          string   `json:"status"`
	SafetyChecklist []string `json:"safetyChecklist"`
	CompletedAt     string   `json:"completedAt,omitempty"`
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
	Executors   []int  `json:"executors"`

	// WorkingHours dikirim sebagai string dari frontend ("0 menit"), disimpan di DB
	WorkingHours string `json:"workingHours"`
	Status       string `json:"status"`

	SafetyChecklist []string `json:"safetyChecklist"`

	// FIX: sql.NullString tidak kompatibel JSON — pakai *string (pointer, bisa nil)
	CompletedAt *string `json:"completedAt,omitempty"`
}

type TakeWorkOrder struct {
	Executors            []int    `json:"executors"`
	SafetyChecklistItems []string `json:"safety_checklist_items"`
	Status               string   `json:"status"`
}

type CompleteWorkOrder struct {
	Status             string `json:"status"`
	CompletedAtDisplay string `json:"completed_at_display"`
}

// UpdateWorkOrderRequest digunakan untuk PATCH /api/workorders/:id
// FIX: struct baru untuk endpoint yang sebelumnya tidak ada
type UpdateWorkOrderRequest struct {
	Executors []int   `json:"executors"`
	Status    *string `json:"status,omitempty"`
}
