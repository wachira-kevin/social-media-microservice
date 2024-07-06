package handlers

import (
	"net/http"
	"notification-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	DB *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

func (h *NotificationHandler) GetNotificationsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	var notifications []models.Notification

	if err := h.DB.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
