package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", RateLimit(1, time.Minute), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := func(ip string) int {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.2:1234"
		req.Header.Set("X-Real-IP", ip)
		router.ServeHTTP(response, req)
		return response.Code
	}

	if got := request("192.0.2.1"); got != http.StatusNoContent {
		t.Fatalf("first request status = %d, want %d", got, http.StatusNoContent)
	}
	if got := request("192.0.2.1"); got != http.StatusTooManyRequests {
		t.Fatalf("second request status = %d, want %d", got, http.StatusTooManyRequests)
	}
	if got := request("192.0.2.2"); got != http.StatusNoContent {
		t.Fatalf("different client status = %d, want %d", got, http.StatusNoContent)
	}
}
