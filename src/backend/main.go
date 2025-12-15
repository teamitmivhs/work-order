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

	dbErr := config.InitDB()
	if dbErr != nil {
		println("Warning: Database connection failed - " + dbErr.Error())
		println("Server akan tetap berjalan tanpa database")
		println("Data akan disimpan di localStorage di frontend")
	} else {
		println("Database connected successfully")
	}
	defer config.CloseDB()

	r := gin.Default()

	setupMiddleware(r)
	setupStaticRoutes(r)
	setupPageRoutes(r)
	setupAPIRoutes(r)

	println("Server starting on port 8080")
	println("Frontend available at: http://localhost:8080")
	println("Summary page: http://localhost:8080/summary")
	println("Kaizen page: http://localhost:8080/kaizen")
	println("TechGuide page: http://localhost:8080/techguide")
	println("")

	if err := r.Run(port); err != nil {
		println("Failed to start server:", err.Error())
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

	r.GET("/techguide", func(c *gin.Context) {
		c.File("../techguide.html")
	})
	r.GET("/src/techguide.html", func(c *gin.Context) {
		c.File("../techguide.html")
	})

	r.GET("/summary", func(c *gin.Context) {
		c.File("../summary.html")
	})
	r.GET("/src/summary.html", func(c *gin.Context) {
		c.File("../summary.html")
	})

	r.GET("/kaizen", func(c *gin.Context) {
		c.File("../kaizen.html")
	})
	r.GET("/src/kaizen.html", func(c *gin.Context) {
		c.File("../kaizen.html")
	})
	// Handle other routes to serve index.html (for SPA support)

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		switch path {
		case "/techguide.html":
			c.File("../techguide.html")
		case "/summary.html":
			c.File("../summary.html")
		case "/kaizen.html":
			c.File("../kaizen.html")
		default:
			c.File("../index.html")
		}
	})
}

func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	routes.RegisterWorkorderRoutes(api)
	routes.RegisterUserRoutes(api)
}
