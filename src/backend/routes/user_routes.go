package routes

import (
	"teamitmivhs/work-order-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", controllers.Register)
	rg.POST("/login", controllers.Login)
}
