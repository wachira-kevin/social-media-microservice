package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user-service/internal/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	// User routes
	userGroup := router.Group("/users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.GET("", userHandler.GetAllUsers)
		userGroup.PUT("/:id", userHandler.UpdateUser)
	}

	return router
}
