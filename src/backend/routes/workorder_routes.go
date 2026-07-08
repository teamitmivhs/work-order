package routes

import (
	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/controllers"
	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/repository"

	"github.com/gin-gonic/gin"
)

// register workorder routes
func RegisterWorkorderRoutes(api *gin.RouterGroup) {
	db := config.GetDB()
	if db == nil {
		println("ERROR: Database connection required for API routes.")
		return
	}

	workOrderRepo := repository.NewWorkOrderRepository(db)
	workOrderCtrl := controllers.NewWorkOrderController(workOrderRepo)

	api.GET("/workorders/track/:code", workOrderCtrl.TrackOrderHandler)
	api.PATCH("/workorders/track/:code/notes", workOrderCtrl.UpdateTrackedOrderNotesHandler)

	//protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		//members
		protected.GET("/members", controllers.GetMembersHandler)
		protected.PATCH("/members/:id/status", controllers.UpdateMemberStatusHandler)
		protected.GET("/shift/day-counter", controllers.GetShiftDayCounterHandler)

		//kaizen
		protected.GET("/kaizen", workOrderCtrl.GetKaizenHandler)

		notifications := protected.Group("/notifications")
		{
			notifications.GET("/vapid-public-key", controllers.GetPushPublicKeyHandler)
			notifications.POST("/subscribe", controllers.SubscribePushNotificationHandler)
			notifications.POST("/unsubscribe", controllers.UnsubscribePushNotificationHandler)
		}

		//workorders
		workorders := protected.Group("/workorders")
		{
			//list
			workorders.GET("", workOrderCtrl.GetTaskListHandler)

			//create
			workorders.POST("", workOrderCtrl.CreateTaskHandler)

			//take
			workorders.POST("/:id/take", workOrderCtrl.TakeOrderHandler)

			//complete
			workorders.PATCH("/:id/complete", workOrderCtrl.CompleteOrderHandler)

			//documentation
			workorders.POST("/:id/documentation", workOrderCtrl.UploadDocumentationPhotoHandler)

			//notes
			workorders.PATCH("/:id/notes", workOrderCtrl.UpdateNotesHandler)

			//update
			workorders.PATCH("/:id", workOrderCtrl.UpdateOrderHandler)

			//delete
			workorders.DELETE("/:id", workOrderCtrl.DeleteOrderHandler)

			//checklist
			workorders.GET("/:id/checklist", workOrderCtrl.GetSafetyChecklistHandler)
			workorders.PUT("/:id/checklist", workOrderCtrl.UpdateSafetyChecklistHandler)
		}
	}
}
