package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
	minNameLength     = 3
	maxNameLength     = 50
)

type RegisterRequest struct {
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	BatchYear string `json:"batchYear"`
	Division  string `json:"division"`
}

// isStrongPassword validates password strength
func isStrongPassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	// Validasi input
	req.Name = strings.TrimSpace(req.Name)
	req.BatchYear = strings.TrimSpace(req.BatchYear)
	division := normalizeStaffDivision(req.Division)
	if len(req.Name) < minNameLength || len(req.Name) > maxNameLength {
		utils.BadRequest(c, "Username must be between 3 and 50 characters")
		return
	}
	if req.BatchYear == "" {
		utils.BadRequest(c, "Angkatan is required")
		return
	}
	if !isValidBatchNumber(req.BatchYear) {
		utils.BadRequest(c, "Angkatan must be a number, for example 13 or 14")
		return
	}
	if division == "" {
		utils.BadRequest(c, "Invalid division. Must be Soundman, Programmer, Maintenance, or Data Analyst")
		return
	}

	if req.Password == "" {
		utils.BadRequest(c, "Password is required")
		return
	}

	if !isStrongPassword(req.Password) {
		utils.BadRequest(c, "Password must be at least 8 characters with uppercase, lowercase, and digits")
		return
	}

	memberRepo := repository.NewMemberRepository()
	existingMember, err := memberRepo.GetMemberByName(req.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.InternalServerError(c, "Failed to check existing account", err)
		return
	}

	if existingMember != nil {
		if existingMember.Password != "" {
			utils.Conflict(c, "Username already exists or is waiting for admin approval")
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password", err)
		return
	}

	if existingMember != nil {
		roleToSet := "Operator"
		if normalizedRole := normalizeStaffRole(existingMember.Role); normalizedRole != "" {
			roleToSet = normalizedRole
		}
		if err := memberRepo.SetMemberPassword(existingMember.ID, string(hashedPassword), roleToSet, division, req.BatchYear); err != nil {
			utils.InternalServerError(c, "Failed to set password", err)
			return
		}
		utils.RespondWithMessage(c, http.StatusCreated, "Registration successful. Your account is waiting for admin approval.", gin.H{
			"member": gin.H{"id": existingMember.ID, "name": existingMember.Name, "accountStatus": existingMember.AccountStatus},
		})
		return
	}

	hasAdmin, err := memberRepo.HasActiveAdmin()
	if err != nil {
		utils.InternalServerError(c, "Failed to check admin bootstrap state", err)
		return
	}

	accountStatus := "pending"
	role := "Operator"
	canHandle := false
	message := "Registration submitted. Please wait for admin approval."
	if !hasAdmin {
		accountStatus = "active"
		role = "Admin"
		canHandle = true
		message = "First admin account created. Please login."
	}

	member := models.Member{
		Name:               req.Name,
		Password:           string(hashedPassword),
		Role:               role,
		Division:           division,
		Status:             "offduty",
		Avatar:             "default-avatar.png",
		AccountStatus:      accountStatus,
		MembershipStatus:   "active",
		BatchYear:          req.BatchYear,
		CanHandleWorkOrder: canHandle,
	}
	if err := memberRepo.CreateMember(&member); err != nil {
		utils.InternalServerError(c, "Failed to create registration request", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusCreated, message, gin.H{
		"member": gin.H{"id": member.ID, "name": member.Name, "accountStatus": member.AccountStatus},
	})
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Password == "" {
		utils.BadRequest(c, "Username and password are required")
		return
	}

	memberRepo := repository.NewMemberRepository()
	member, err := memberRepo.GetMemberByName(req.Name)
	if err != nil {
		// Don't reveal if user exists for security
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	if member == nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	if member.AccountStatus != "" && member.AccountStatus != "active" {
		utils.Forbidden(c, "Your account is not active yet. Please wait for admin approval.")
		return
	}
	if member.MembershipStatus == "alumni" || member.MembershipStatus == "inactive" {
		utils.Forbidden(c, "Your account is no longer active as staff.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(member.ID, member.Name, member.Role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token", err)
		return
	}

	// Don't send password to frontend
	member.Password = ""
	utils.RespondWithMessage(c, http.StatusOK, "Login successful", gin.H{
		"token":  token,
		"member": member,
	})
}

// GetProfile returns current user profile
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User information not found")
		return
	}

	id, ok := userID.(int)
	if !ok {
		utils.Unauthorized(c, "Invalid user information")
		return
	}

	memberRepo := repository.NewMemberRepository()
	member, err := memberRepo.GetMemberByID(id)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	member.Password = ""
	utils.RespondSuccess(c, http.StatusOK, member)
}

// ChangePasswordHandler: PATCH /api/profile/password
func ChangePasswordHandler(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "User not found")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		utils.BadRequest(c, "Current password and new password are required")
		return
	}
	if !isStrongPassword(req.NewPassword) {
		utils.BadRequest(c, "Password must be at least 8 characters with uppercase, lowercase, and digits")
		return
	}
	if req.CurrentPassword == req.NewPassword {
		utils.BadRequest(c, "New password must be different from current password")
		return
	}

	memberRepo := repository.NewMemberRepository()
	member, err := memberRepo.GetMemberByID(userID)
	if err != nil || member == nil {
		utils.Unauthorized(c, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.CurrentPassword)); err != nil {
		utils.Unauthorized(c, "Current password is incorrect")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password", err)
		return
	}
	if err := memberRepo.UpdateMemberPassword(userID, string(hashedPassword)); err != nil {
		utils.InternalServerError(c, "Failed to update password", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Password updated successfully", nil)
}

// UploadAvatarHandler: POST /api/profile/avatar
// Upload foto profil user — simpan ke static/public/ dan update DB
func UploadAvatarHandler(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "User not found")
		return
	}

	// Ambil file dari form
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.BadRequest(c, "No file uploaded")
		return
	}
	defer file.Close()

	// Validasi ukuran (max 5MB). Foto dari HP sering >2MB, sementara nginx
	// juga sudah dinaikkan limit-nya agar request tidak ditolak sebelum sampai backend.
	if header.Size > 5*1024*1024 {
		utils.BadRequest(c, "File too large. Maximum 5MB")
		return
	}

	// Validasi extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		utils.BadRequest(c, "Invalid file type. Allowed: jpg, jpeg, png, webp")
		return
	}

	// Generate nama file unik — pakai userID + timestamp agar tidak bentrok
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)

	// Pastikan folder static/public ada
	uploadDir := utils.PublicUploadDir()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.InternalServerError(c, "Failed to create upload directory", err)
		return
	}

	// Simpan file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		utils.InternalServerError(c, "Failed to save file", err)
		return
	}

	// Update kolom Avatar di database
	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.UpdateMemberAvatar(userID, filename); err != nil {
		utils.InternalServerError(c, "Failed to update avatar in database", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Avatar updated successfully", gin.H{
		"avatar": filename,
		"url":    "/static/public/" + filename,
	})
}

// DeleteAvatarHandler: DELETE /api/profile/avatar
// Reset avatar ke default
func DeleteAvatarHandler(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "User not found")
		return
	}

	memberRepo := repository.NewMemberRepository()

	// Ambil avatar lama untuk dihapus dari disk
	member, err := memberRepo.GetMemberByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	// Hapus file lama jika bukan default
	if member.Avatar != "" && member.Avatar != "no avatar" &&
		!strings.HasPrefix(member.Avatar, "default") {
		oldPath := filepath.Join(utils.PublicUploadDir(), member.Avatar)
		os.Remove(oldPath) // silent fail — file mungkin sudah tidak ada
	}

	// Reset ke default di DB
	if err := memberRepo.UpdateMemberAvatar(userID, "no avatar"); err != nil {
		utils.InternalServerError(c, "Failed to reset avatar", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Avatar removed successfully", gin.H{
		"avatar": "no avatar",
	})
}

// UpdateStatusHandler: POST /api/status
// Update status member (standby, onjob, support, nextshift, offduty)
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func UpdateStatusHandler(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "User not found")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	// Trim dan lowercase status
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))

	// Validasi status (validation juga ada di repository, tapi cek dulu di sini)
	validStatuses := map[string]bool{
		"standby": true, "onjob": true, "support": true,
		"nextshift": true, "offduty": true,
	}
	if !validStatuses[req.Status] {
		utils.BadRequest(c, "Invalid status. Must be: standby, onjob, support, nextshift, or offduty")
		return
	}

	memberRepo := repository.NewMemberRepository()

	// Update status di database
	if err := memberRepo.UpdateMemberStatus(userID, req.Status); err != nil {
		utils.InternalServerError(c, "Failed to update status", err)
		return
	}

	// Ambil data member terbaru setelah update
	member, err := memberRepo.GetMemberByID(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch updated member data", err)
		return
	}

	member.Password = ""
	utils.RespondWithMessage(c, http.StatusOK, "Status updated successfully", gin.H{
		"member": member,
		"status": member.Status,
	})
}

type ApproveMemberRequest struct {
	Role               string `json:"role"`
	Division           string `json:"division"`
	BatchYear          string `json:"batchYear"`
	CanHandleWorkOrder *bool  `json:"canHandleWorkOrder"`
}

type GraduateBatchRequest struct {
	BatchYear      string `json:"batchYear" binding:"required"`
	GraduationYear int    `json:"graduationYear" binding:"required"`
}

func GetAdminMembersHandler(c *gin.Context) {
	filter := strings.ToLower(strings.TrimSpace(c.Query("status")))
	memberRepo := repository.NewMemberRepository()
	members, err := memberRepo.GetAdminMembers(filter)
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve staff list", err)
		return
	}

	for i := range members {
		members[i].Password = ""
	}
	utils.RespondSuccess(c, http.StatusOK, members)
}

func ApproveMemberHandler(c *gin.Context) {
	memberID, ok := parseMemberIDParam(c)
	if !ok {
		return
	}

	approverID, ok := middleware.GetUserIDFromContext(c)
	if !ok || approverID == 0 {
		utils.Unauthorized(c, "Admin user not found")
		return
	}

	var req ApproveMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	role := normalizeStaffRole(req.Role)
	if role == "" {
		utils.BadRequest(c, "Invalid role. Must be Operator or Admin")
		return
	}

	division := normalizeStaffDivision(req.Division)
	if division == "" {
		utils.BadRequest(c, "Invalid division. Must be Soundman, Programmer, Maintenance, or Data Analyst")
		return
	}
	batchYear := strings.TrimSpace(req.BatchYear)
	if !isValidBatchNumber(batchYear) {
		utils.BadRequest(c, "Angkatan must be a number, for example 13 or 14")
		return
	}

	canHandle := true
	if req.CanHandleWorkOrder != nil {
		canHandle = *req.CanHandleWorkOrder
	}

	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.ApproveMember(memberID, role, division, batchYear, canHandle, approverID); err != nil {
		utils.InternalServerError(c, "Failed to approve staff", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Staff approved successfully", gin.H{"id": memberID})
}

func RejectMemberHandler(c *gin.Context) {
	memberID, ok := parseMemberIDParam(c)
	if !ok {
		return
	}
	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.RejectMember(memberID); err != nil {
		utils.InternalServerError(c, "Failed to reject staff", err)
		return
	}
	utils.RespondWithMessage(c, http.StatusOK, "Staff rejected successfully", gin.H{"id": memberID})
}

func DisableMemberHandler(c *gin.Context) {
	memberID, ok := parseMemberIDParam(c)
	if !ok {
		return
	}
	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.DisableMember(memberID); err != nil {
		utils.InternalServerError(c, "Failed to disable staff", err)
		return
	}
	utils.RespondWithMessage(c, http.StatusOK, "Staff disabled successfully", gin.H{"id": memberID})
}

func MarkMemberAlumniHandler(c *gin.Context) {
	memberID, ok := parseMemberIDParam(c)
	if !ok {
		return
	}

	var req struct {
		GraduationYear int `json:"graduationYear" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}
	if req.GraduationYear < 2000 || req.GraduationYear > 2100 {
		utils.BadRequest(c, "Graduation year is invalid")
		return
	}

	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.MarkMemberAsAlumni(memberID, req.GraduationYear); err != nil {
		utils.InternalServerError(c, "Failed to mark staff as alumni", err)
		return
	}
	utils.RespondWithMessage(c, http.StatusOK, "Staff moved to alumni successfully", gin.H{"id": memberID})
}

func GraduateBatchHandler(c *gin.Context) {
	var req GraduateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	req.BatchYear = strings.TrimSpace(req.BatchYear)
	if req.BatchYear == "" {
		utils.BadRequest(c, "Batch year is required")
		return
	}
	if !isValidBatchNumber(req.BatchYear) {
		utils.BadRequest(c, "Angkatan must be a number, for example 13 or 14")
		return
	}
	if req.GraduationYear < 2000 || req.GraduationYear > 2100 {
		utils.BadRequest(c, "Graduation year is invalid")
		return
	}

	memberRepo := repository.NewMemberRepository()
	affected, err := memberRepo.GraduateBatch(req.BatchYear, req.GraduationYear)
	if err != nil {
		utils.InternalServerError(c, "Failed to graduate batch", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Batch graduated successfully", gin.H{"affected": affected})
}

// ChangeRoleHandler: PATCH /api/admin/members/:id/role
// Allows an admin to promote Operator → Admin or demote Admin → Operator.
// An admin cannot change their own role (guard against self-lockout).
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func ChangeRoleHandler(c *gin.Context) {
	memberID, ok := parseMemberIDParam(c)
	if !ok {
		return
	}

	callerID, ok := middleware.GetUserIDFromContext(c)
	if !ok || callerID == 0 {
		utils.Unauthorized(c, "Admin user not found")
		return
	}

	if callerID == memberID {
		utils.BadRequest(c, "You cannot change your own role")
		return
	}

	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	role := normalizeStaffRole(req.Role)
	if role == "" {
		utils.BadRequest(c, "Invalid role. Must be Operator or Admin")
		return
	}

	memberRepo := repository.NewMemberRepository()
	if err := memberRepo.ChangeRole(memberID, role); err != nil {
		utils.InternalServerError(c, "Failed to change role", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Role updated successfully", gin.H{
		"id":   memberID,
		"role": role,
	})
}

func parseMemberIDParam(c *gin.Context) (int, bool) {
	memberID, err := strconv.Atoi(c.Param("id"))
	if err != nil || memberID <= 0 {
		utils.BadRequest(c, "Invalid member ID")
		return 0, false
	}
	return memberID, true
}

func normalizeStaffRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "", "operator":
		return "Operator"
	case "admin":
		return "Admin"
	default:
		return ""
	}
}

func normalizeStaffDivision(division string) string {
	switch strings.ToLower(strings.TrimSpace(division)) {
	case "soundman":
		return "Soundman"
	case "", "programmer":
		return "Programmer"
	case "maintenance":
		return "Maintenance"
	case "data analyst", "data-analyst", "data_analyst":
		return "Data Analyst"
	default:
		return ""
	}
}

func isValidBatchNumber(batch string) bool {
	batch = strings.TrimSpace(batch)
	if batch == "" || len(batch) > 2 {
		return false
	}
	value, err := strconv.Atoi(batch)
	return err == nil && value >= 1 && value <= 99
}
