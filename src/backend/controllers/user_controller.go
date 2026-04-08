package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	// Validasi input
	member.Name = strings.TrimSpace(member.Name)
	if len(member.Name) < minNameLength || len(member.Name) > maxNameLength {
		utils.BadRequest(c, "Username must be between 3 and 50 characters")
		return
	}

	if member.Password == "" {
		utils.BadRequest(c, "Password is required")
		return
	}

	if !isStrongPassword(member.Password) {
		utils.BadRequest(c, "Password must be at least 8 characters with uppercase, lowercase, and digits")
		return
	}

	// Check apakah nama sudah ada di DB
	// PENTING: Hanya member yang sudah terdaftar di database yang bisa melakukan registration
	memberRepo := repository.NewMemberRepository()
	existingMember, err := memberRepo.GetMemberByName(member.Name)

	// Jika ada error atau nama tidak ditemukan di database, tolak registration
	if err != nil || existingMember == nil {
		utils.BadRequest(c, "Username not found in system. You are not authorized to register. Please contact IT administrator.")
		return
	}

	// Member ada di database — cek apakah passwordnya masih kosong (belum pernah register)
	// Ini terjadi karena member di-seed dari SQL dengan password kosong
	if existingMember.Password != "" {
		// Sudah punya password — tolak, tidak boleh overwrite akun aktif
		utils.Conflict(c, "Username already exists")
		return
	}

	// Password masih kosong — ini member yang belum pernah register
	// Set password mereka dan kembalikan token
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(member.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password", err)
		return
	}

	// Tentukan role: pertahankan role yang ada dari DB (programmer, maintenance, dll)
	// Upgrade ke 'Operator' agar bisa akses protected endpoints
	roleToSet := "Operator"
	if existingMember.Role != "" {
		roleToSet = existingMember.Role
	}

	if err := memberRepo.SetMemberPassword(existingMember.ID, string(hashedPassword), roleToSet); err != nil {
		utils.InternalServerError(c, "Failed to set password", err)
		return
	}

	token, err := utils.GenerateToken(existingMember.ID, existingMember.Name, roleToSet)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusCreated, "Registration successful", gin.H{
		"token":  token,
		"member": gin.H{"id": existingMember.ID, "name": existingMember.Name, "role": roleToSet, "status": existingMember.Status},
	})
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	// Validasi ukuran (max 2MB)
	if header.Size > 2*1024*1024 {
		utils.BadRequest(c, "File too large. Maximum 2MB")
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
	uploadDir := "/static/public"
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
		oldPath := filepath.Join("/static/public", member.Avatar)
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
