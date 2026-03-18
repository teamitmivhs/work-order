package routes

import (
	"net/http"
	"sync"
	"time"

	"teamitmivhs/work-order-backend/controllers"
	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

// ─── Simple In-Process Rate Limiter ──────────────────────────────────────────
// FIX: endpoint /register dan /login sebelumnya tidak ada rate limiting sama sekali
// sehingga rawan brute-force dan spam registrasi.
//
// Implementasi ini memakai token bucket sederhana per IP.
// Untuk production dengan multiple instance, ganti dengan Redis-based rate limiter
// (misalnya github.com/ulule/limiter dengan Redis store).

type rateBucket struct {
	tokens    int
	lastRefil time.Time
}

type ipRateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*rateBucket
	maxToken int           // max request per window
	window   time.Duration // window refill
}

func newIPRateLimiter(maxToken int, window time.Duration) *ipRateLimiter {
	rl := &ipRateLimiter{
		buckets:  make(map[string]*rateBucket),
		maxToken: maxToken,
		window:   window,
	}
	// Bersihkan bucket lama setiap 10 menit agar tidak memory leak
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			rl.mu.Lock()
			cutoff := time.Now().Add(-window * 2)
			for ip, b := range rl.buckets {
				if b.lastRefil.Before(cutoff) {
					delete(rl.buckets, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

// allow mengembalikan true jika request dari IP ini masih dalam batas
func (rl *ipRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.buckets[ip]
	if !exists {
		rl.buckets[ip] = &rateBucket{tokens: rl.maxToken - 1, lastRefil: time.Now()}
		return true
	}

	// Refill token jika window sudah lewat
	if time.Since(b.lastRefil) >= rl.window {
		b.tokens = rl.maxToken
		b.lastRefil = time.Now()
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Limiter untuk auth endpoint: 10 request per menit per IP
var authLimiter = newIPRateLimiter(10, time.Minute)

// rateLimitMiddleware adalah Gin middleware yang memakai authLimiter
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !authLimiter.allow(ip) {
			utils.RespondError(c, http.StatusTooManyRequests,
				"Too many requests. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ─── Routes ──────────────────────────────────────────────────────────────────

func RegisterUserRoutes(rg *gin.RouterGroup) {
	// Public routes dengan rate limiting
	// FIX: sebelumnya tidak ada rate limiting — rawan brute-force
	public := rg.Group("")
	public.Use(rateLimitMiddleware())
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// Protected routes (butuh JWT)
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)

		// FIX: tambah endpoint logout
		// JWT stateless tidak bisa di-invalidate di server secara native,
		// tapi endpoint ini memberi tempat yang jelas untuk:
		// 1. Membersihkan token di sisi client (response menginstruksikan ini)
		// 2. Di masa depan: blacklist token di Redis sebelum expiry
		protected.POST("/logout", logoutHandler)
	}
}

// logoutHandler menangani request logout.
// Saat ini menginstruksikan client untuk menghapus token lokal.
// Jika di masa depan ada Redis, token bisa diblacklist di sini.
func logoutHandler(c *gin.Context) {
	// Ambil info user dari context untuk logging audit
	userName, _ := middleware.GetUserNameFromContext(c)
	if userName != "" {
		utils.LogInfo("User '%s' logged out", userName)
	}

	utils.RespondWithMessage(c, http.StatusOK,
		"Logged out successfully. Please remove your token on the client side.", nil)
}
