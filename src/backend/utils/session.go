package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const SessionCookieName = "work_order_session"

func SetSessionCookie(c *gin.Context, token string, maxAge int) {
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionCookieName, token, maxAge, "/", "", secure, true)
}
