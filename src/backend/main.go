package main

import (
	"time"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	port = ":8080"
)

func main() {
	println("🚀 Starting Work Order Server...")

	// Initialize database connection (optional - server runs without DB)
	dbErr := config.InitDB()
	if dbErr != nil {
		println("⚠️  Warning: Database connection failed - " + dbErr.Error())
		println("📝 Server akan tetap berjalan tanpa database")
		println("💾 Data akan disimpan di localStorage di frontend")
	} else {
		println("✅ Database connected successfully")
	}
	defer config.CloseDB()

	// Create Gin router
	r := gin.Default()

	setupMiddleware(r)
	setupStaticRoutes(r)
	setupPageRoutes(r)
	setupAPIRoutes(r)

	// Start server
	println("🌐 Server starting on port 8080")
	println("📱 Frontend available at: http://localhost:8080")
	println("📊 Summary page: http://localhost:8080/summary")
	println("🎯 Kaizen page: http://localhost:8080/kaizen")
	println("📋 TechGuide page: http://localhost:8080/techguide")
	println("")

	if err := r.Run(port); err != nil {
		println("❌ Failed to start server:", err.Error())
	}
}

func setupMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func setupStaticRoutes(r *gin.Engine) {
	r.Static("/static", "../static")
	r.Static("/src/static", "../static")
}

func setupPageRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.File("../index.html")
	})

	// TechGuide route
	r.GET("/techguide", func(c *gin.Context) {
		c.File("../techguide.html")
	})
	r.GET("/src/techguide.html", func(c *gin.Context) {
		c.File("../techguide.html")
	})

	// Summary route
	r.GET("/summary", func(c *gin.Context) {
		c.File("../summary.html")
	})
	r.GET("/src/summary.html", func(c *gin.Context) {
		c.File("../summary.html")
	})

	// Kaizen route
	r.GET("/kaizen", func(c *gin.Context) {
		c.File("../kaizen.html")
	})
	r.GET("/src/kaizen.html", func(c *gin.Context) {
		c.File("../kaizen.html")
	})

	// Fallback to index for any unmatched route
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/techguide.html" {
			c.File("../techguide.html")
		} else if path == "/summary.html" {
			c.File("../summary.html")
		} else if path == "/kaizen.html" {
			c.File("../kaizen.html")
		} else {
			c.File("../index.html")
		}
	})
}

func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	routes.RegisterWorkorderRoutes(api)
}
