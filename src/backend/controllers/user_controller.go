package controllers

import (
	"net/http"
	"strings"

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

// isStrongPassword memvalidasi kekuatan password:
// minimal 8 karakter, mengandung huruf besar, huruf kecil, dan angka
func isStrongPassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

// Register membuat akun member baru
func Register(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

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

	// Cek duplikasi username
	existingMember, err := memberRepo.GetMemberByName(member.Name)
	if err == nil && existingMember != nil {
		utils.Conflict(c, "Username already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(member.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password", err)
		return
	}
	member.Password = string(hashedPassword)

	// Default role dan status untuk member baru
	member.Role = "Operator"
	member.Status = "standby"

	if err := memberRepo.CreateMember(&member); err != nil {
		utils.InternalServerError(c, "Failed to create member", err)
		return
	}

	token, err := utils.GenerateToken(member.ID, member.Name, member.Role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusCreated, "Registration successful", gin.H{
		"token": token,
		"member": gin.H{
			"id":     member.ID,
			"name":   member.Name,
			"role":   member.Role,
			"status": member.Status,
		},
	})
}

type LoginRequest struct {
	Name     string `json:"name"     binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login memvalidasi kredensial dan mengembalikan JWT token
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
	if err != nil || member == nil {
		// Jangan ungkapkan apakah username ada atau tidak (security best practice)
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(member.ID, member.Name, member.Role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token", err)
		return
	}

	// Jangan kirim password ke frontend
	member.Password = ""
	utils.RespondWithMessage(c, http.StatusOK, "Login successful", gin.H{
		"token":  token,
		"member": member,
	})
}

// GetProfile mengembalikan profil user yang sedang login
func GetProfile(c *gin.Context) {
	// FIX: pakai signature baru GetUserIDFromContext yang return (int, bool)
	// Sebelumnya return 0 jika tidak ada, yang bisa lolos pengecekan ceroboh
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		utils.Unauthorized(c, "User information not found")
		return
	}

	memberRepo := repository.NewMemberRepository()
	member, err := memberRepo.GetMemberByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	member.Password = ""
	utils.RespondSuccess(c, http.StatusOK, member)
}
