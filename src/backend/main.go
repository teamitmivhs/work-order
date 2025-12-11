package main

import (
	"time"

	"teamitmivhs/work-order-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", "http://127.0.0.1:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve frontend static files from the parent `src/` directory so
	// requests to `/` and static assets are served by this server.
	// Register static handlers BEFORE API routes to avoid wildcard conflicts
	// with existing `/api` prefixes.
	// Serve static assets under /static to avoid conflicts with /api routes.
	// This maps requests like /static/assets/script.js -> src/assets/script.js
	r.Static("/static", "../")

	// Root serves index.html
	r.GET("/", func(c *gin.Context) {
		c.File("../index.html")
	})

	// TechGuide page
	r.GET("/techguide", func(c *gin.Context) {
		c.File("../techguide.html")
	})

	// Fallback for unknown paths (optional SPA support)
	r.NoRoute(func(c *gin.Context) {
		c.File("../index.html")
	})

	api := r.Group("/api")
	routes.RegisterWorkorderRoutes(api)

	r.Run(":8080")
}
