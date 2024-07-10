package routers

import (
	"github.com/gin-gonic/gin"
	"post-service/internal/handlers"
)

func SetupRouter(postHandler *handlers.PostHandler) *gin.Engine {
	r := gin.Default()

	postGroup := r.Group("/posts")
	{
		postGroup.POST("/", postHandler.CreatePost)
		postGroup.GET("/", postHandler.GetAllPosts)
		postGroup.PUT("/:id", postHandler.UpdatePost)
		postGroup.GET("/:id", postHandler.GetPostByID)
		postGroup.GET("/user/:user_id", postHandler.GetPostsByUserID)
		postGroup.POST("/:post_id/like/:user_id", postHandler.LikePost)
		postGroup.POST("/:post_id/comment", postHandler.CommentOnPost)
	}

	return r
}
