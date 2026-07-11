package controllers

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/services"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

func generateTrackingCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[n.Int64()]
	}
	return "WO-" + string(code), nil
}

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

// GetTaskListHandler menangani request GET /api/workorders
// Protected endpoint — Guest tidak boleh melihat dashboard internal.
func (ctrl *WorkOrderController) GetTaskListHandler(c *gin.Context) {
	var tasks []models.WorkOrder
	var err error

	userRole, _ := middleware.GetUserRoleFromContext(c)

	if userRole == "Admin" || userRole == "Operator" {
		tasks, err = ctrl.Repo.GetAllTasks()
	} else if userRole == "Guest" {
		utils.Forbidden(c, "Guest tidak dapat melihat daftar work order")
		return
	} else {
		utils.Forbidden(c, "Work order access requires an internal account")
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

	if role, _ := middleware.GetUserRoleFromContext(c); role != "Admin" && role != "Guest" {
		utils.Forbidden(c, "Only admins and guests can create work orders")
		return
	}

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
	if req.TrackingCode == "" {
		trackingCode, err := generateTrackingCode()
		if err != nil {
			utils.InternalServerError(c, "Failed to generate tracking code", err)
			return
		}
		req.TrackingCode = trackingCode
	}

	newID, err := ctrl.Repo.CreateTask(req)
	if err != nil {
		utils.InternalServerError(c, "Failed to save task", err)
		return
	}

	go services.NotifyNewWorkOrder(newID, req.Device, req.Location, req.Priority)

	utils.RespondWithMessage(c, http.StatusCreated, "Work order created successfully", gin.H{
		"id":           newID,
		"trackingCode": req.TrackingCode,
	})
}

func (ctrl *WorkOrderController) TrackOrderHandler(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		utils.BadRequest(c, "Tracking code is required")
		return
	}

	order, err := ctrl.Repo.GetTaskByTrackingCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFound(c, "Work order tidak ditemukan")
			return
		}
		utils.InternalServerError(c, "Failed to track work order", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"id":                 order.ID,
		"trackingCode":       order.TrackingCode,
		"status":             order.Status,
		"priority":           order.Priority,
		"time":               order.Time,
		"requester":          order.Requester,
		"location":           order.Location,
		"device":             order.Device,
		"problem":            order.Problem,
		"executors":          order.Executors,
		"completedAt":        order.CompletedAt,
		"workingHours":       order.WorkingHours,
		"notes":              order.Notes,
		"adminNotes":         order.AdminNotes,
		"documentationPhoto": order.DocumentationPhoto,
	})
}

func (ctrl *WorkOrderController) UpdateTrackedOrderNotesHandler(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		utils.BadRequest(c, "Tracking code is required")
		return
	}

	var req models.UpdateNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	notes := strings.TrimSpace(strings.NewReplacer("<", "", ">", "").Replace(req.Notes))
	if len(notes) > 1000 {
		utils.BadRequest(c, "Catatan maksimal 1000 karakter")
		return
	}

	order, err := ctrl.Repo.GetTaskByTrackingCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFound(c, "Work order tidak ditemukan")
			return
		}
		utils.InternalServerError(c, "Failed to track work order", err)
		return
	}

	if err := ctrl.Repo.UpdateOrderNoteText(int64(order.ID), notes); err != nil {
		utils.InternalServerError(c, "Failed to save notes", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Notes saved successfully", gin.H{
		"id":           order.ID,
		"trackingCode": order.TrackingCode,
		"notes":        notes,
	})
}

// TakeOrderHandler: POST /api/workorders/{id}/take
// Hanya member yang di-assign yang bisa take order
func (ctrl *WorkOrderController) TakeOrderHandler(c *gin.Context) {
	role, _ := middleware.GetUserRoleFromContext(c)
	if role == "Guest" {
		utils.Forbidden(c, "Guest hanya bisa membuat work order")
		return
	}

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

	if role != "Admin" {
		userID, ok := middleware.GetUserIDFromContext(c)
		if !ok || userID == 0 {
			utils.Unauthorized(c, "User not found")
			return
		}
		if err := ctrl.Repo.JoinPendingOrder(orderID, userID); err != nil {
			utils.InternalServerError(c, "Failed to join work order", err)
			return
		}
		operatorName := fmt.Sprintf("Operator #%d", userID)
		if member, err := ctrl.MemberRepo.GetMemberByID(userID); err == nil && member.Name != "" {
			operatorName = member.Name
		}
		go services.NotifyWorkOrderAdmins(
			"Approval work order menunggu",
			fmt.Sprintf("%s mengambil work order #%d. Approve untuk mulai progress.", operatorName, orderID),
			orderID,
			nil,
		)
		utils.RespondWithMessage(c, http.StatusOK, "Operator added. Waiting for admin approval.", gin.H{"id": orderID})
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

	// TakeOrder tidak perlu cek IsMemberAssigned karena executor baru diisi
	// di dalam TakeOrder itu sendiri — belum ada sebelum take dilakukan
	err = ctrl.Repo.TakeOrder(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to take order", err)
		return
	}

	go services.NotifyWorkOrderUsers(
		req.Executors,
		"Work order mulai progress",
		fmt.Sprintf("Work order #%d sudah di-approve dan mulai progress.", orderID),
		orderID,
		nil,
	)

	utils.RespondWithMessage(c, http.StatusOK, "Order taken successfully", gin.H{"id": orderID})
}

// CompleteOrderHandler: PATCH /api/workorders/{id}/complete
// Validasi: hanya assigned member, safety checklist fulfilled
func (ctrl *WorkOrderController) CompleteOrderHandler(c *gin.Context) {
	// Guest tidak boleh complete order
	if role, _ := middleware.GetUserRoleFromContext(c); role == "Guest" {
		utils.Forbidden(c, "Guest hanya bisa membuat work order")
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	userRole, _ := middleware.GetUserRoleFromContext(c)

	var req models.CompleteWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.Status != "completed" {
		utils.BadRequest(c, "Status must be 'completed'")
		return
	}

	// Cek assignment hanya untuk Operator — Admin bisa complete order apapun
	if userRole != "Admin" {
		userID, _ := middleware.GetUserIDFromContext(c)
		isAssigned, err := ctrl.MemberRepo.IsMemberAssigned(orderID, userID)
		if err != nil {
			utils.InternalServerError(c, "Failed to check assignment", err)
			return
		}
		if !isAssigned {
			utils.Forbidden(c, "You are not assigned to this work order")
			return
		}
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

	hasPhoto, err := ctrl.Repo.HasDocumentationPhoto(orderID)
	if err != nil {
		utils.InternalServerError(c, "Failed to check documentation photo", err)
		return
	}
	if !hasPhoto {
		utils.BadRequest(c, "Documentation photo is required before finishing the work order")
		return
	}

	err = ctrl.Repo.CompleteOrder(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to complete order", err)
		return
	}

	go services.NotifyWorkOrderBroadcast(
		"Work order selesai",
		fmt.Sprintf("Work order #%d sudah ditandai selesai.", orderID),
		orderID,
		nil,
	)

	utils.RespondWithMessage(c, http.StatusOK, "Order completed successfully", gin.H{"id": orderID})
}

func (ctrl *WorkOrderController) RejectOrderHandler(c *gin.Context) {
	if role, _ := middleware.GetUserRoleFromContext(c); role != "Admin" {
		utils.Forbidden(c, "Only admins can reject work orders")
		return
	}

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	var req models.RejectWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}
	reason := strings.TrimSpace(strings.NewReplacer("<", "", ">", "").Replace(req.Reason))
	if len(reason) < 3 {
		utils.BadRequest(c, "Reject reason is required")
		return
	}
	if len(reason) > 1000 {
		utils.BadRequest(c, "Reject reason max 1000 characters")
		return
	}

	if err := ctrl.Repo.RejectOrder(orderID, reason); err != nil {
		utils.BadRequest(c, "Failed to reject order", err.Error())
		return
	}

	go services.NotifyWorkOrderBroadcast(
		"Work order ditolak",
		fmt.Sprintf("Work order #%d ditolak. Alasan: %s", orderID, reason),
		orderID,
		nil,
	)

	utils.RespondWithMessage(c, http.StatusOK, "Order rejected successfully", gin.H{"id": orderID})
}

// UploadDocumentationPhotoHandler: POST /api/workorders/{id}/documentation
// Upload foto bukti pekerjaan. Hanya admin atau operator assigned yang boleh upload.
func (ctrl *WorkOrderController) UploadDocumentationPhotoHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	userRole, _ := middleware.GetUserRoleFromContext(c)
	userID, _ := middleware.GetUserIDFromContext(c)
	if userRole == "Guest" {
		utils.Forbidden(c, "Guest cannot upload work order documentation")
		return
	}
	if userRole != "Admin" {
		isAssigned, err := ctrl.MemberRepo.IsMemberAssigned(orderID, userID)
		if err != nil {
			utils.InternalServerError(c, "Failed to check assignment", err)
			return
		}
		if !isAssigned {
			utils.Forbidden(c, "You are not assigned to this work order")
			return
		}
	}

	file, err := c.FormFile("photo")
	if err != nil {
		utils.BadRequest(c, "No documentation photo uploaded")
		return
	}
	if file.Size > 15*1024*1024 {
		utils.BadRequest(c, "File too large. Maximum 15MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		utils.BadRequest(c, "Invalid file type. Allowed: jpg, jpeg, png, webp")
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		utils.BadRequest(c, "Invalid documentation photo")
		return
	}
	defer openedFile.Close()
	if err := utils.ValidateImageUpload(openedFile, map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}); err != nil {
		utils.BadRequest(c, "Invalid image file")
		return
	}

	uploadDir := filepath.Join(utils.PublicUploadDir(), "workorder-docs")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.InternalServerError(c, "Failed to create documentation upload directory", err)
		return
	}

	filename := fmt.Sprintf("workorder_%d_%d%s", orderID, time.Now().Unix(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.InternalServerError(c, "Failed to save documentation photo", err)
		return
	}

	relativeFilename := "workorder-docs/" + filename
	if err := ctrl.Repo.UpdateOrderDocumentationPhoto(orderID, relativeFilename); err != nil {
		utils.InternalServerError(c, "Failed to update documentation photo", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Documentation photo uploaded successfully", gin.H{
		"documentationPhoto": relativeFilename,
		"url":                "/static/public/" + relativeFilename,
	})
}

// DeleteOrderHandler: DELETE /api/workorders/{id}
// Hanya admin yang bisa delete
func (ctrl *WorkOrderController) DeleteOrderHandler(c *gin.Context) {
	// Guest tidak boleh delete order
	if role, _ := middleware.GetUserRoleFromContext(c); role != "Admin" {
		utils.Forbidden(c, "Only admins can delete work orders")
		return
	}

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
	if role, _ := middleware.GetUserRoleFromContext(c); role != "Admin" {
		utils.Forbidden(c, "Only admins can update safety checklist")
		return
	}

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
	if _, err := repository.RunShiftDayRollover(); err != nil {
		utils.InternalServerError(c, "Failed to run shift day counter", err)
		return
	}

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

func GetShiftDayCounterHandler(c *gin.Context) {
	snapshot, err := repository.RunShiftDayRollover()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve shift day counter", err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, snapshot)
}

// UpdateOrderHandler: PATCH /api/workorders/{id}
// Update executor list untuk order yang masih pending
func (ctrl *WorkOrderController) UpdateOrderHandler(c *gin.Context) {
	if role, _ := middleware.GetUserRoleFromContext(c); role != "Admin" {
		utils.Forbidden(c, "Only admins can update work orders")
		return
	}

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

	err = ctrl.Repo.UpdateOrderExecutors(orderID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to update order", err)
		return
	}

	go services.NotifyWorkOrderUsers(
		req.Executors,
		"Ditambahkan ke work order",
		fmt.Sprintf("Kamu ditambahkan ke work order #%d.", orderID),
		orderID,
		nil,
	)

	utils.RespondWithMessage(c, http.StatusOK, "Order updated successfully", gin.H{"id": orderID})
}

// UpdateNotesHandler: PATCH /api/workorders/{id}/notes
// Simpan catatan evaluasi untuk work order yang sudah selesai
func (ctrl *WorkOrderController) UpdateNotesHandler(c *gin.Context) {
	role, _ := middleware.GetUserRoleFromContext(c)
	if role != "Admin" {
		utils.Forbidden(c, "Only admins can update notes and ratings")
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid Order ID")
		return
	}

	var req models.UpdateNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		utils.BadRequest(c, "Rating must be between 1 and 5")
		return
	}
	if req.NotesQuality != nil && (*req.NotesQuality < 1 || *req.NotesQuality > 5) {
		utils.BadRequest(c, "Notes quality must be between 1 and 5")
		return
	}

	adminNotes := strings.TrimSpace(strings.NewReplacer("<", "", ">", "").Replace(req.AdminNotes))
	if len(adminNotes) > 1000 {
		utils.BadRequest(c, "Catatan admin maksimal 1000 karakter")
		return
	}

	if err := ctrl.Repo.UpdateOrderNotes(orderID, adminNotes, req.Rating, req.NotesQuality); err != nil {
		utils.InternalServerError(c, "Failed to save notes", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Notes saved successfully", gin.H{"id": orderID})
}

// UpdateMemberStatusHandler: PATCH /api/members/{id}/status
// Update status member (standby, onjob, nextshift, offduty).
// Admin dan Data Analyst boleh update semua member; operator hanya boleh update status dirinya sendiri.
func UpdateMemberStatusHandler(c *gin.Context) {
	memberIDStr := c.Param("id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid Member ID")
		return
	}

	userRole, _ := middleware.GetUserRoleFromContext(c)
	userID, _ := middleware.GetUserIDFromContext(c)
	memberRepo := repository.NewMemberRepository()
	canManageShift := userRole == "Admin"
	if !canManageShift && userID != 0 {
		currentMember, err := memberRepo.GetMemberByID(userID)
		if err == nil && strings.EqualFold(strings.TrimSpace(currentMember.Division), "Data Analyst") {
			canManageShift = true
		}
	}
	if !canManageShift && userID != memberID {
		utils.Forbidden(c, "Only admins or Data Analyst can update other member statuses")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	validStatuses := map[string]bool{
		"standby": true, "onjob": true,
		"nextshift": true, "offduty": true,
	}
	if !validStatuses[req.Status] {
		utils.BadRequest(c, "Invalid status. Must be: standby, onjob, nextshift, or offduty")
		return
	}

	if err := memberRepo.UpdateMemberStatus(memberID, req.Status); err != nil {
		utils.InternalServerError(c, "Failed to update member status", err)
		return
	}

	member, err := memberRepo.GetMemberByID(memberID)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch updated member data", err)
		return
	}
	member.Password = ""

	utils.RespondWithMessage(c, http.StatusOK, "Member status updated successfully", gin.H{
		"id":     memberID,
		"status": member.Status,
		"member": member,
	})
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
