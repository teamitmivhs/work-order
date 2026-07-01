package models

type WorkOrder struct {
	ID                 int      `json:"id"`
	Priority           string   `json:"priority"`
	Time               string   `json:"time"`
	StartedAt          string   `json:"startedAt,omitempty"`
	ProgressSeconds    *int     `json:"progressSeconds,omitempty"`
	TrackingCode       string   `json:"trackingCode,omitempty"`
	Requester          string   `json:"requester"`
	Location           string   `json:"location"`
	Device             string   `json:"device"`
	Problem            string   `json:"problem"`
	Executors          []int    `json:"executors"`
	WorkingHours       *int     `json:"workingHours,omitempty"`
	Status             string   `json:"status"`
	SafetyChecklist    []string `json:"safetyChecklist"`
	CompletedAt        string   `json:"completedAt,omitempty"`
	Notes              string   `json:"notes,omitempty"`
	Rating             *int     `json:"rating,omitempty"`
	NotesQuality       *int     `json:"notesQuality,omitempty"`
	DocumentationPhoto string   `json:"documentationPhoto,omitempty"`
}

type Member struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Password           string `json:"password,omitempty"`
	Role               string `json:"role"`
	Division           string `json:"division,omitempty"`
	Status             string `json:"status"`
	Avatar             string `json:"avatar"`
	AccountStatus      string `json:"accountStatus,omitempty"`
	MembershipStatus   string `json:"membershipStatus,omitempty"`
	BatchYear          string `json:"batchYear,omitempty"`
	GraduationYear     *int   `json:"graduationYear,omitempty"`
	CanHandleWorkOrder bool   `json:"canHandleWorkOrder"`
	RegisteredAt       string `json:"registeredAt,omitempty"`
	ApprovedAt         string `json:"approvedAt,omitempty"`
	ApprovedBy         *int   `json:"approvedBy,omitempty"`
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
	ID           int    `json:"id"`
	Priority     string `json:"priority"`
	TimeDisplay  string `json:"time_display"`
	TimeSort     string `json:"time_sort"`
	TrackingCode string `json:"trackingCode,omitempty"`
	Requester    string `json:"requester"`
	Location     string `json:"location"`
	Device       string `json:"device"`
	Problem      string `json:"problem"`

	// Executors harus diisi dari tabel relasi task_executors
	Executors []int `json:"executors"`

	// WorkingHours bisa NULL di DB
	WorkingHours string `json:"workingHours"`
	Status       string `json:"status"`

	// SafetyChecklist harus diisi dari tabel relasi task_safety_checklists
	SafetyChecklist []string `json:"safetyChecklist"`

	// CompletedAt bisa NULL di DB
	CompletedAt *string `json:"completedAt,omitempty"`
}

// UpdateWorkOrderRequest: payload untuk PATCH /api/workorders/:id
// Dipakai untuk update executor list dan status order
type UpdateWorkOrderRequest struct {
	Executors []int   `json:"executors"`
	Status    *string `json:"status,omitempty"`
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

// UpdateNotesRequest: payload untuk PATCH /api/workorders/:id/notes
type UpdateNotesRequest struct {
	Notes        string `json:"notes"`
	Rating       *int   `json:"rating"`
	NotesQuality *int   `json:"notesQuality"`
}
