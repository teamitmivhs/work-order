package routes

import (
	"teamitmivhs/work-order-backend/controllers"
	"teamitmivhs/work-order-backend/middleware"

	"net/http"

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
		protected.POST("/logout", logoutHandler)
		protected.POST("/profile/avatar", controllers.UploadAvatarHandler)
		protected.DELETE("/profile/avatar", controllers.DeleteAvatarHandler)
	}
}

// logoutHandler — stateless JWT: instruksikan client hapus token
func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Logged out successfully. Please remove your token on the client side.",
	})
}
