package main

import (
	"time"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/routes"
	"teamitmivhs/work-order-backend/services"

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
		if snapshot, err := repository.RunShiftDayRollover(); err != nil {
			println("Warning: shift day counter failed - " + err.Error())
		} else if snapshot.MovedToStandby > 0 {
			println("Shift rollover moved members to standby:", snapshot.MovedToStandby)
		}
		services.InitPushNotifications()
	}
	defer config.CloseDB()

	r := gin.Default()

	//Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	setupMiddleware(r)
	setupStaticRoutes(r)
	setupPageRoutes(r)
	setupAPIRoutes(r)

	println("Server starting on port 8080")
	println("Frontend available at: http://localhost:8080")
	println("Summary page: http://localhost:8080/summary")
	println("Kaizen page: http://localhost:8080/kaizen")
	println("Shift page: http://localhost:8080/shift")
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
	// Auth pages
	r.GET("/login.html", func(c *gin.Context) {
		c.File("../login.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("../login.html")
	})

	r.GET("/register.html", func(c *gin.Context) {
		c.File("../register.html")
	})
	r.GET("/register", func(c *gin.Context) {
		c.File("../register.html")
	})

	r.GET("/guest", func(c *gin.Context) {
		c.File("../guest.html")
	})
	r.GET("/guest.html", func(c *gin.Context) {
		c.File("../guest.html")
	})

	r.GET("/staff", func(c *gin.Context) {
		c.File("../staff.html")
	})
	r.GET("/staff.html", func(c *gin.Context) {
		c.File("../staff.html")
	})

	r.GET("/shift", func(c *gin.Context) {
		c.File("../shift.html")
	})
	r.GET("/shift.html", func(c *gin.Context) {
		c.File("../shift.html")
	})

	// Dashboard
	r.GET("/", func(c *gin.Context) {
		c.File("../index.html")
	})

	// TechGuide page
	r.GET("/techguide", func(c *gin.Context) {
		c.File("../techguide.html")
	})
	r.GET("/src/techguide.html", func(c *gin.Context) {
		c.File("../techguide.html")
	})
	r.GET("/techguide.html", func(c *gin.Context) {
		c.File("../techguide.html")
	})

	// Summary page
	r.GET("/summary", func(c *gin.Context) {
		c.File("../summary.html")
	})
	r.GET("/src/summary.html", func(c *gin.Context) {
		c.File("../summary.html")
	})
	r.GET("/summary.html", func(c *gin.Context) {
		c.File("../summary.html")
	})

	// Kaizen page
	r.GET("/kaizen", func(c *gin.Context) {
		c.File("../kaizen.html")
	})
	r.GET("/src/kaizen.html", func(c *gin.Context) {
		c.File("../kaizen.html")
	})
	r.GET("/kaizen.html", func(c *gin.Context) {
		c.File("../kaizen.html")
	})

	// Fallback: serve index.html for any unmatched route (SPA behavior)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		switch path {
		case "/login.html":
			c.File("../login.html")
		case "/register.html":
			c.File("../register.html")
		case "/guest.html":
			c.File("../guest.html")
		case "/staff.html":
			c.File("../staff.html")
		case "/shift.html":
			c.File("../shift.html")
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
