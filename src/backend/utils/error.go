package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// LogError logs error dengan context
func LogError(message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	} else {
		log.Printf("[ERROR] %s", message)
	}
}

// LogInfo logs informational messages
func LogInfo(message string) {
	log.Printf("[INFO] %s", message)
}

// RespondError mengirim error response dengan status code yang sesuai
func RespondError(c *gin.Context, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	if len(details) > 0 {
		response.Details = details[0]
	}
	c.JSON(statusCode, response)
}

// RespondSuccess mengirim success response
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	if statusCode == http.StatusNoContent {
		c.Status(statusCode)
		return
	}
	c.JSON(statusCode, data)
}

// RespondWithMessage mengirim response dengan message
func RespondWithMessage(c *gin.Context, statusCode int, message string, data interface{}) {
	response := map[string]interface{}{
		"message": message,
	}
	if data != nil {
		response["data"] = data
	}
	c.JSON(statusCode, response)
}

// Common error responses
func Unauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message)
}

func BadRequest(c *gin.Context, message string, details ...string) {
	if len(details) > 0 {
		RespondError(c, http.StatusBadRequest, message, details[0])
	} else {
		RespondError(c, http.StatusBadRequest, message)
	}
}

func NotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string, err error) {
	LogError(message, err)
	RespondError(c, http.StatusInternalServerError, message)
}

func Conflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, message)
}
