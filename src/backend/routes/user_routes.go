package routes

import (
	"teamitmivhs/work-order-backend/controllers"
	"teamitmivhs/work-order-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	// Public routes (no auth required)
	rg.POST("/register", controllers.Register)
	rg.POST("/login", controllers.Login)

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)
	}
}
