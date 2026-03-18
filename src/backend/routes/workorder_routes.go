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

	// ─── Public endpoints (tidak butuh auth) ──────────────────────────────────
	// FIX: GET /members tetap public agar frontend bisa load member list
	api.GET("/members", controllers.GetMembersHandler)

	// FIX: GET /workorders dipindah ke public agar dashboard bisa load tanpa login
	// (frontend tidak mengirim JWT token saat GET workorders)
	// Jika di masa depan ingin diprotect, pindahkan ke grup protected di bawah
	api.GET("/workorders", workOrderCtrl.GetTaskListHandler)

	// ─── Protected endpoints (butuh JWT) ──────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Kaizen metrics
		protected.GET("/kaizen", workOrderCtrl.GetKaizenHandler)

		// FIX: PATCH /members/:id/status — dibutuhkan frontend untuk update status member
		protected.PATCH("/members/:id/status", controllers.UpdateMemberStatusHandler)

		workorders := protected.Group("/workorders")
		{
			// Buat work order baru
			workorders.POST("", workOrderCtrl.CreateTaskHandler)

			// Take order (jadikan progress)
			workorders.POST("/:id/take", workOrderCtrl.TakeOrderHandler)

			// Tandai selesai
			workorders.PATCH("/:id/complete", workOrderCtrl.CompleteOrderHandler)

			// FIX: PATCH /:id — untuk addWorkerToOrder & syncRemoveMemberFromOrders dari frontend
			// Sebelumnya endpoint ini tidak ada sama sekali
			workorders.PATCH("/:id", workOrderCtrl.UpdateOrderHandler)

			// Hapus order (Admin only)
			workorders.DELETE("/:id", middleware.AdminMiddleware(), workOrderCtrl.DeleteOrderHandler)

			// Safety checklist
			workorders.GET("/:id/checklist", workOrderCtrl.GetSafetyChecklistHandler)
			workorders.PUT("/:id/checklist", workOrderCtrl.UpdateSafetyChecklistHandler)
		}
	}
}
