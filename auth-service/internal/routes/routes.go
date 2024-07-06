package routes

import (
	"auth-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout/:sessionId", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	return r
}
