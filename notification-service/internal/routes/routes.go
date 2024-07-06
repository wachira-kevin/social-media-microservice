package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"notification-service/internal/handlers"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	notificationHandler := handlers.NewNotificationHandler(db)

	router := gin.Default()

	// User routes
	notificationsGroup := router.Group("/notifications")
	{
		notificationsGroup.GET("/events/:id", handlers.SSEHandler)
		notificationsGroup.GET("/:userID", notificationHandler.GetNotificationsByUserID)
	}

	return router
}
