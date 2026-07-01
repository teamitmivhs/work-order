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
		protected.PATCH("/profile/password", controllers.ChangePasswordHandler)
		protected.POST("/profile/avatar", controllers.UploadAvatarHandler)
		protected.DELETE("/profile/avatar", controllers.DeleteAvatarHandler)
		protected.POST("/status", controllers.UpdateStatusHandler)

		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/members", controllers.GetAdminMembersHandler)
			admin.PATCH("/members/:id/approve", controllers.ApproveMemberHandler)
			admin.PATCH("/members/:id/reject", controllers.RejectMemberHandler)
			admin.PATCH("/members/:id/disable", controllers.DisableMemberHandler)
			admin.PATCH("/members/:id/alumni", controllers.MarkMemberAlumniHandler)
			admin.PATCH("/members/:id/role", controllers.ChangeRoleHandler)
			admin.POST("/members/graduate", controllers.GraduateBatchHandler)
		}
	}
}

func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Logged out successfully. Please remove your token on the client side.",
	})
}
