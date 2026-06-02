package controllers

import (
	"net/http"
	"strconv"
	"sync"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

type WorkOrderController struct {
	Repo       repository.WorkOrderRepository
	MemberRepo repository.MemberRepository
}

var (
	mu         sync.Mutex
	workOrders = []models.WorkOrder{}
	nextID     = 1
)

func NewWorkOrderController(repo repository.WorkOrderRepository) *WorkOrderController {
	return &WorkOrderController{
		Repo:       repo,
		MemberRepo: repository.NewMemberRepository(),
	}
}

// GetTaskListHandler menangani request GET /api/workorders
// Filter berdasarkan role user
func (ctrl *WorkOrderController) GetTaskListHandler(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	userRole := middleware.GetUserRoleFromContext(c)

	var tasks []models.WorkOrder
	var err error

	// Filter berdasarkan role
	switch userRole {
	case "Admin":
		// Admin melihat semua orders
		tasks, err = ctrl.Repo.GetAllTasks()
	case "Operator":
		// Operator hanya melihat orders yang mereka assigned
		tasks, err = ctrl.Repo.GetTasksByExecutor(userID)
	default:
		utils.Forbidden(c, "Invalid role")
		return
	}

	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve tasks", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, tasks)
}

// CreateTaskHandler menangani request POST /api/workorders
func (ctrl *WorkOrderController) CreateTaskHandler(c *gin.Context) {
	var req models.WorkOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	// Validasi input
	if req.Priority == "" {
		utils.BadRequest(c, "Priority is required")
		return
	}
	if req.Requester == "" {
		utils.BadRequest(c, "Requester is required")
		return
	}
	if req.Location == "" {
		utils.BadRequest(c, "Location is required")
		return
	}
	if req.Device == "" {
		utils.BadRequest(c, "Device is required")
		return
	}
	if req.Problem == "" {
		utils.BadRequest(c, "Problem description is required")
		return
	}

	// Validasi priority
	validPriorities := []string{"low", "medium", "high", "urgent"}
	isValidPriority := false
	for _, p := range validPriorities {
		if req.Priority == p {
			isValidPriority = true
			break
		}
	}
	if !isValidPriority {
		utils.BadRequest(c, "Invalid priority. Must be: low, medium, high, or urgent")
		return
	}

	// Set default status
	if req.Status == "" {
		req.Status = "pending"
	}

	newID, err := ctrl.Repo.CreateTask(req)
	if err != nil {
		utils.InternalServerError(c, "Failed to save task", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusCreated, "Work order created successfully", gin.H{"id": newID})
}

// TakeOrderHandler: POST /api/workorders/{id}/take
// Hanya member yang di-assign yang bisa take order
func (ctrl *WorkOrderController) TakeOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	userID := middleware.GetUserIDFromContext(c)

	var req models.TakeWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.Status != "progress" {
		utils.BadRequest(c, "Status must be 'progress' to take order")
		return
	}

	// Check if member is assigned to this order
	isAssigned, err := ctrl.MemberRepo.IsMemberAssigned(orderID, userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to check assignment", err)
		return
	}

	if !isAssigned {
		utils.Forbidden(c, "You are not assigned to this work order")
		return
	}

	err = ctrl.Repo.TakeOrder(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to take order", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Order taken successfully", gin.H{"id": orderID})
}

// CompleteOrderHandler: PATCH /api/workorders/{id}/complete
// Validasi: hanya assigned member, safety checklist fulfilled
func (ctrl *WorkOrderController) CompleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	userID := middleware.GetUserIDFromContext(c)

	var req models.CompleteWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.Status != "completed" {
		utils.BadRequest(c, "Status must be 'completed'")
		return
	}

	// Check if member is assigned to this order
	isAssigned, err := ctrl.MemberRepo.IsMemberAssigned(orderID, userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to check assignment", err)
		return
	}

	if !isAssigned {
		utils.Forbidden(c, "You are not assigned to this work order")
		return
	}

	// Validasi safety checklist
	checklistFulfilled, err := ctrl.Repo.IsSafetyChecklistFulfilled(orderID)
	if err != nil {
		utils.InternalServerError(c, "Failed to check safety checklist", err)
		return
	}

	if !checklistFulfilled {
		utils.BadRequest(c, "Safety checklist must be completed before finishing the work order")
		return
	}

	err = ctrl.Repo.CompleteOrder(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to complete order", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Order completed successfully", gin.H{"id": orderID})
}

// DeleteOrderHandler: DELETE /api/workorders/{id}
// Hanya admin yang bisa delete
func (ctrl *WorkOrderController) DeleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	err = ctrl.Repo.DeleteOrder(orderID)
	if err != nil {
		utils.InternalServerError(c, "Failed to delete order", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Order deleted successfully", gin.H{"id": orderID})
}

// GetSafetyChecklistHandler: GET /api/workorders/{id}/checklist
func (ctrl *WorkOrderController) GetSafetyChecklistHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	checklist, err := ctrl.Repo.GetSafetyChecklist(orderID)
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve safety checklist", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"checklist": checklist})
}

// UpdateSafetyChecklistHandler: PUT /api/workorders/{id}/checklist
func (ctrl *WorkOrderController) UpdateSafetyChecklistHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	var req struct {
		ChecklistItems []string `json:"checklist_items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if len(req.ChecklistItems) == 0 {
		utils.BadRequest(c, "Checklist items cannot be empty")
		return
	}

	err = ctrl.Repo.UpdateSafetyChecklist(orderID, req.ChecklistItems)
	if err != nil {
		utils.InternalServerError(c, "Failed to update safety checklist", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Safety checklist updated successfully", nil)
}

// GetKaizenHandler: GET /api/kaizen
// Return kaizen/performance metrics
func (ctrl *WorkOrderController) GetKaizenHandler(c *gin.Context) {
	metrics, err := ctrl.Repo.GetKaizenMetrics()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve kaizen metrics", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, metrics)
}

// GetMembersHandler: GET /api/members
func GetMembersHandler(c *gin.Context) {
	memberRepo := repository.NewMemberRepository()
	members, err := memberRepo.GetAllMembers()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve members", err)
		return
	}

	// Don't send passwords
	for i := range members {
		members[i].Password = ""
	}

	utils.RespondSuccess(c, http.StatusOK, members)
}

// OLD IN-MEMORY FUNCTIONS - DEPRECATED (KEEPING FOR REFERENCE BUT NOT USED)
// These have been replaced with database-driven implementations above

/*
func CreateWorkOrder(c *gin.Context) { ... }
func UpdateWorkOrder(c *gin.Context) { ... }
func DeleteWorkOrder(c *gin.Context) { ... }
func GetSummary(c *gin.Context) { ... }
func GetKaizen(c *gin.Context) { ... }
*/
