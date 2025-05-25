package controller

import (
	"net/http"

	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService services.NotificationService
}

func NewNotificationController(notificationService services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: notificationService}
}

func (c *NotificationController) SendMessage(ctx *gin.Context) {
	notification := &model.Notification{}

	if err := ctx.ShouldBindJSON(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the SendMessage method to send the notification
	sentNotification, err := c.notificationService.SendMessage(notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notification": sentNotification})
}
