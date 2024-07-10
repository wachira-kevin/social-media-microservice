package handlers

import (
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"net/http"
	"post-service/internal/models"
	"post-service/internal/publishers"
	"post-service/internal/services"
	"strconv"
)

type PostHandler struct {
	PostService *services.PostService
	conn        *amqp.Connection
}

func NewPostHandler(postService *services.PostService, conn *amqp.Connection) *PostHandler {
	return &PostHandler{
		PostService: postService,
		conn:        conn,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var newPost models.CreatePost
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create post to db
	post := &models.Post{
		UserID:  newPost.UserID,
		Content: newPost.Content,
	}
	if err := h.PostService.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// publish message to create new post notification
	event := &models.CreatePostNotificationEvent{
		PostID:  post.ID,
		UserID:  post.UserID,
		Content: post.Content,
	}
	err := publishers.PublishNewPostMessage(event, h.conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postToUpdate, err := h.PostService.GetPostByID(uint(id))
	postToUpdate.Content = post.Content
	if err := h.PostService.UpdatePost(postToUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, postToUpdate)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	post, err := h.PostService.GetPostByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	posts, err := h.PostService.GetPostsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) LikePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	post, err := h.PostService.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.PostService.LikePost(uint(postID), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add like event notification
	event := &models.CreateLikeNotificationEvent{
		PostID:  post.ID,
		UserID:  post.UserID,
		LikerID: uint(userID),
	}
	err = publishers.PublishNewLikeMessage(event, h.conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

func (h *PostHandler) CommentOnPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var newComment models.CreateComment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := h.PostService.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// save comment
	comment := &models.Comment{
		PostID:  uint(postID),
		UserID:  newComment.UserID,
		Content: newComment.Content,
	}
	if err := h.PostService.CommentOnPost(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// send new Comment
	event := &models.CreateCommentNotificationEvent{
		PostID:      uint(postID),
		UserID:      post.UserID,
		CommenterID: newComment.UserID,
		Comment:     newComment.Content,
	}
	err = publishers.PublishNewCommentMessage(event, h.conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, comment)
}
