package main

import (
	"time"

	"teamitmivhs/work-order-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	port = ":8080"
)

func main() {
	r := gin.Default()

	setupMiddleware(r)
	setupStaticRoutes(r)
	setupPageRoutes(r)
	setupAPIRoutes(r)

	r.Run(port)
}

func setupMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", "http://127.0.0.1:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
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
		c.File("src/index.html")
	})

	// TechGuide route
	r.GET("/techguide", func(c *gin.Context) {
		c.File("src/techguide.html")
	})
	r.GET("/src/techguide.html", func(c *gin.Context) {
		c.File("src/techguide.html")
	})

	// Summary route
	r.GET("/summary", func(c *gin.Context) {
		c.File("src/summary.html")
	})
	r.GET("/src/summary.html", func(c *gin.Context) {
		c.File("src/summary.html")
	})

	// Fallback to index for any unmatched route
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/techguide.html" {
			c.File("src/techguide.html")
		} else if path == "/summary.html" {
			c.File("src/summary.html")
		} else {
			c.File("src/index.html")
		}
	})
}

func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	routes.RegisterWorkorderRoutes(api)
}
