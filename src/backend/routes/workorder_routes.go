package routes

import (
	"teamitmivhs/work-order-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterWorkorderRoutes(rg *gin.RouterGroup) {
	rg.GET("/workorders", controllers.GetWorkOrders)
	rg.POST("/workorders", controllers.CreateWorkOrder)
	rg.PUT("/workorders/:id", controllers.UpdateWorkOrder)
	rg.PATCH("/workorders/:id", controllers.UpdateWorkOrder)
	rg.DELETE("/workorders/:id", controllers.DeleteWorkOrder)

	rg.GET("/summary", controllers.GetSummary)
	rg.GET("/kaizen", controllers.GetKaizen)
	rg.GET("/members", controllers.GetMembersHandler)
}
