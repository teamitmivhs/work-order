package routes

import (
	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/controllers"
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

	// Daftarkan endpoint untuk mengambil semua member
	api.GET("/members", controllers.GetMembersHandler)
	api.GET("/summary", controllers.GetSummary)
	api.GET("/kaizen", controllers.GetKaizen)

	workorders := api.Group("/workorders")
	{
		// CREATE
		workorders.POST("", workOrderCtrl.CreateTaskHandler)

		// READ
		workorders.GET("", workOrderCtrl.GetTaskListHandler) // Ambil semua order

		// OPERASI KRUSIAL
		workorders.POST("/:id/take", workOrderCtrl.TakeOrderHandler)          // Ambil Order (Start Progress)
		workorders.PATCH("/:id/complete", workOrderCtrl.CompleteOrderHandler) // Tandai Selesai
		workorders.DELETE("/:id", workOrderCtrl.DeleteOrderHandler)           // Hapus Order
	}
}
