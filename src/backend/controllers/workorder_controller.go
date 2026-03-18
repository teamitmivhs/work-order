package controllers

import (
	"net/http"
	"strconv"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

// FIX: hapus variabel global in-memory yang tidak dipakai (mu, workOrders, nextID)
// dan import "sync" yang tidak perlu

type WorkOrderController struct {
	Repo       repository.WorkOrderRepository
	MemberRepo repository.MemberRepository
}

func NewWorkOrderController(repo repository.WorkOrderRepository) *WorkOrderController {
	return &WorkOrderController{
		Repo:       repo,
		MemberRepo: repository.NewMemberRepository(),
	}
}

// GetTaskListHandler: GET /api/workorders
// Admin melihat semua, Operator hanya melihat yang di-assign ke mereka
func (ctrl *WorkOrderController) GetTaskListHandler(c *gin.Context) {
	// FIX: pakai signature baru (T, bool) untuk deteksi "tidak ada user" secara eksplisit
	userRole, roleOK := middleware.GetUserRoleFromContext(c)
	if !roleOK {
		// Jika tidak ada role (request tanpa auth / public access), tampilkan semua
		// sesuai keputusan di routes bahwa GET /workorders adalah public
		tasks, err := ctrl.Repo.GetAllTasks()
		if err != nil {
			utils.InternalServerError(c, "Failed to retrieve tasks", err)
			return
		}
		utils.RespondSuccess(c, http.StatusOK, tasks)
		return
	}

	userID, _ := middleware.GetUserIDFromContext(c)

	var tasks []models.WorkOrder
	var err error

	switch userRole {
	case "Admin":
		tasks, err = ctrl.Repo.GetAllTasks()
	case "Operator":
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

// CreateTaskHandler: POST /api/workorders
func (ctrl *WorkOrderController) CreateTaskHandler(c *gin.Context) {
	var req models.WorkOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

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

	validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "urgent": true}
	if !validPriorities[req.Priority] {
		utils.BadRequest(c, "Invalid priority. Must be: low, medium, high, or urgent")
		return
	}

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

// TakeOrderHandler: POST /api/workorders/:id/take
// FIX: hapus pengecekan IsMemberAssigned sebelum take — logika sebelumnya terbalik.
// Saat take, user justru belum di-assign. Pengecekan assignment yang benar ada di CompleteOrderHandler.
// Di sini cukup validasi bahwa executor yang dikirim tidak sedang onjob (dilakukan di repo layer).
func (ctrl *WorkOrderController) TakeOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	var req models.TakeWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.Status != "progress" {
		utils.BadRequest(c, "Status must be 'progress' to take order")
		return
	}

	if len(req.Executors) == 0 {
		utils.BadRequest(c, "At least one executor is required")
		return
	}

	err = ctrl.Repo.TakeOrder(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to take order", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Order taken successfully", gin.H{"id": orderID})
}

// CompleteOrderHandler: PATCH /api/workorders/:id/complete
// FIX: pengecekan IsMemberAssigned tetap di sini — hanya member yang assigned boleh complete
func (ctrl *WorkOrderController) CompleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	// FIX: pakai signature baru (int, bool)
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		utils.Unauthorized(c, "User information not found")
		return
	}

	var req models.CompleteWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.Status != "completed" {
		utils.BadRequest(c, "Status must be 'completed'")
		return
	}

	// Hanya member yang di-assign ke order ini yang boleh complete
	isAssigned, err := ctrl.MemberRepo.IsMemberAssigned(orderID, userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to check assignment", err)
		return
	}
	if !isAssigned {
		utils.Forbidden(c, "You are not assigned to this work order")
		return
	}

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

// DeleteOrderHandler: DELETE /api/workorders/:id
// Hanya Admin (dijaga oleh AdminMiddleware di routes)
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

// UpdateOrderHandler: PATCH /api/workorders/:id
// FIX: endpoint baru — dibutuhkan frontend untuk addWorkerToOrder & syncRemoveMemberFromOrders
func (ctrl *WorkOrderController) UpdateOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	var req models.UpdateWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := ctrl.Repo.UpdateOrderExecutors(orderID, req); err != nil {
		utils.InternalServerError(c, "Failed to update order", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Order updated successfully", gin.H{"id": orderID})
}

// GetSafetyChecklistHandler: GET /api/workorders/:id/checklist
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

// UpdateSafetyChecklistHandler: PUT /api/workorders/:id/checklist
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
func (ctrl *WorkOrderController) GetKaizenHandler(c *gin.Context) {
	metrics, err := ctrl.Repo.GetKaizenMetrics()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve kaizen metrics", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, metrics)
}

// GetMembersHandler: GET /api/members
// Public endpoint — tidak butuh auth
func GetMembersHandler(c *gin.Context) {
	memberRepo := repository.NewMemberRepository()
	members, err := memberRepo.GetAllMembers()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve members", err)
		return
	}

	for i := range members {
		members[i].Password = ""
	}

	utils.RespondSuccess(c, http.StatusOK, members)
}

// UpdateMemberStatusHandler: PATCH /api/members/:id/status
// FIX: endpoint baru — frontend memanggil ini saat admin ubah status member
func UpdateMemberStatusHandler(c *gin.Context) {
	memberIDStr := c.Param("id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid member ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.UpdateMemberStatus(memberID, req.Status); err != nil {
		utils.InternalServerError(c, "Failed to update member status", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Member status updated successfully", gin.H{"id": memberID})
}
