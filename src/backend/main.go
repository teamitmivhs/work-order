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

	r.Static("/static", "../")

	r.GET("/", func(c *gin.Context) {
		c.File("../index.html")
	})

	r.GET("/techguide", func(c *gin.Context) {
		c.File("../techguide.html")
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("../index.html")
	})

	api := r.Group("/api")
	routes.RegisterWorkorderRoutes(api)

	r.Run(":8080")
}
