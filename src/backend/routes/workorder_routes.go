package routes

import (
	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/controllers"
	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/repository"

	"github.com/gin-gonic/gin"
)

// RegisterWorkorderRoutes mendaftarkan semua endpoint yang berhubungan dengan Work Order
func RegisterWorkorderRoutes(api *gin.RouterGroup) {
	db := config.GetDB()
	if db == nil {
		println("ERROR: Database connection required for API routes.")
		return
	}

	workOrderRepo := repository.NewWorkOrderRepository(db)
	workOrderCtrl := controllers.NewWorkOrderController(workOrderRepo)

	// Public endpoint (no auth required)
	api.GET("/members", controllers.GetMembersHandler)

	// Protected endpoints (auth required)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Kaizen metrics - accessible to all authenticated users
		protected.GET("/kaizen", workOrderCtrl.GetKaizenHandler)

		workorders := protected.Group("/workorders")
		{
			// Create work order
			workorders.POST("", workOrderCtrl.CreateTaskHandler)

			// Get work orders (filtered by role)
			workorders.GET("", workOrderCtrl.GetTaskListHandler)

			// Work order operations
			workorders.POST("/:id/take", workOrderCtrl.TakeOrderHandler)
			workorders.PATCH("/:id/complete", workOrderCtrl.CompleteOrderHandler)
			workorders.DELETE("/:id", middleware.AdminMiddleware(), workOrderCtrl.DeleteOrderHandler)

			// Safety checklist endpoints
			workorders.GET("/:id/checklist", workOrderCtrl.GetSafetyChecklistHandler)
			workorders.PUT("/:id/checklist", workOrderCtrl.UpdateSafetyChecklistHandler)
		}
	}
}
