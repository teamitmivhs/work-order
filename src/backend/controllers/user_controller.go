package controllers

import (
	"net/http"
	"strings"
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

	memberRepo := repository.NewMemberRepository()

	// Whitelist check: nama harus sudah ada di tabel members (36 member IT MIVHS)
	existingMember, err := memberRepo.GetMemberByName(member.Name)
	if err != nil || existingMember == nil {
		utils.BadRequest(c, "Name not found. Only registered IT MIVHS members can create an account.")
		return
	}

	// Cek apakah sudah pernah register (password sudah diset)
	if existingMember.Password != "" {
		utils.Conflict(c, "Account already registered. Please login instead.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(member.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password", err)
		return
	}

	// Role: jika masih role lama (job title), ubah ke Operator
	// Admin harus diset manual di DB setelah register
	role := existingMember.Role
	if role == "" || role == "programmer" || role == "maintenance" ||
		role == "data analyst" || role == "soundman" {
		role = "Operator"
	}

	// Update password member yang sudah ada (bukan INSERT baru)
	if err := memberRepo.SetMemberPassword(existingMember.ID, string(hashedPassword), role); err != nil {
		utils.InternalServerError(c, "Failed to register account", err)
		return
	}

	token, err := utils.GenerateToken(existingMember.ID, existingMember.Name, role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusCreated, "Registration successful", gin.H{
		"token":  token,
		"member": gin.H{"id": existingMember.ID, "name": existingMember.Name, "role": role, "status": existingMember.Status},
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
