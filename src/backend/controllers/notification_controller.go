package controllers

import (
	"net/http"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/services"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

type pushSubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func GetPushPublicKeyHandler(c *gin.Context) {
	key := services.PushPublicKey()
	if key == "" {
		utils.InternalServerError(c, "Push notifications are not configured", nil)
		return
	}
	utils.RespondSuccess(c, http.StatusOK, gin.H{"publicKey": key})
}

func SubscribePushNotificationHandler(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		utils.Unauthorized(c, "User not found")
		return
	}

	var req pushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid push subscription", err.Error())
		return
	}

	if err := services.SavePushSubscription(userID, req.Endpoint, req.Keys.P256DH, req.Keys.Auth, c.Request.UserAgent()); err != nil {
		utils.BadRequest(c, "Failed to save push subscription", err.Error())
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Push subscription saved", nil)
}

func UnsubscribePushNotificationHandler(c *gin.Context) {
	var req pushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid push subscription", err.Error())
		return
	}

	if err := services.DeletePushSubscription(req.Endpoint); err != nil {
		utils.InternalServerError(c, "Failed to delete push subscription", err)
		return
	}

	utils.RespondWithMessage(c, http.StatusOK, "Push subscription deleted", nil)
}
