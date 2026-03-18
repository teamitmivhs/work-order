package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── Response Structs ────────────────────────────────────────────────────────

// ErrorResponse adalah format standar untuk semua error response.
// Selalu memiliki field: code, message, dan opsional details.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse adalah format standar untuk semua success response.
// FIX: sebelumnya RespondWithMessage dan RespondSuccess menggunakan format berbeda
// (satu punya field "code", satu tidak) sehingga frontend harus handle dua format.
// Sekarang semua response sukses konsisten: { code, message, data }.
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ─── Logging ─────────────────────────────────────────────────────────────────

// LogError mencatat error ke log dengan prefix [ERROR].
// FIX: sebelumnya log.Printf dipanggil langsung di banyak tempat — tidak konsisten.
// Semua logging kini melewati fungsi ini agar mudah diganti library di masa depan.
func LogError(message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	} else {
		log.Printf("[ERROR] %s", message)
	}
}

// LogWarn mencatat warning ke log dengan prefix [WARN].
func LogWarn(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[WARN] "+message, args...)
	} else {
		log.Printf("[WARN] %s", message)
	}
}

// LogInfo mencatat informasi ke log dengan prefix [INFO].
func LogInfo(message string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[INFO] "+message, args...)
	} else {
		log.Printf("[INFO] %s", message)
	}
}

// ─── Core Responders ─────────────────────────────────────────────────────────

// RespondError mengirim error response dengan format standar.
//
// PENTING: fungsi ini TIDAK memanggil c.Abort().
// Caller yang memanggil dari dalam middleware WAJIB memanggil c.Abort() sendiri
// agar handler chain tidak dilanjutkan. Contoh yang benar:
//
//	utils.Unauthorized(c, "...")
//	c.Abort()
//	return
//
// Controller handler tidak perlu c.Abort() karena cukup dengan return.
func RespondError(c *gin.Context, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	if len(details) > 0 && details[0] != "" {
		response.Details = details[0]
	}
	c.JSON(statusCode, response)
}

// RespondSuccess mengirim success response.
// FIX: format sekarang konsisten menggunakan SuccessResponse { code, message, data }
// Sebelumnya format berbeda dengan RespondWithMessage sehingga frontend
// harus handle dua struktur JSON yang berbeda.
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	if statusCode == http.StatusNoContent {
		c.Status(statusCode)
		return
	}
	c.JSON(statusCode, SuccessResponse{
		Code: statusCode,
		Data: data,
	})
}

// RespondWithMessage mengirim success response dengan pesan dan data opsional.
// FIX: format sekarang konsisten dengan RespondSuccess — keduanya pakai SuccessResponse.
func RespondWithMessage(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// ─── Shorthand Error Helpers ─────────────────────────────────────────────────
// Semua helper di bawah adalah wrapper tipis di atas RespondError.
// Middleware yang memakainya WAJIB memanggil c.Abort() setelahnya (lihat catatan di atas).

func Unauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message)
}

func BadRequest(c *gin.Context, message string, details ...string) {
	RespondError(c, http.StatusBadRequest, message, details...)
}

func NotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

// InternalServerError mencatat error ke log lalu mengirim response 500.
// Detail error TIDAK dikirim ke client untuk menghindari kebocoran informasi internal.
func InternalServerError(c *gin.Context, message string, err error) {
	LogError(message, err)
	RespondError(c, http.StatusInternalServerError, message)
}

func Conflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, message)
}
