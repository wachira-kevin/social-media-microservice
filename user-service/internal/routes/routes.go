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
		userGroup.PUT("/:id/settings", userHandler.UpdateUserSettings)
		userGroup.POST("/:follower_id/follow/:followee_id", userHandler.FollowUser)
		userGroup.GET("/:id/followers", userHandler.GetFollowers)
		userGroup.GET("/:id/following", userHandler.GetFollowing)
	}

	return router
}
