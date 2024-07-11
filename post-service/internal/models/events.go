package models

type CreatePostNotificationEvent struct {
	// Define fields for post notification event
	PostID  uint   `json:"post_id"`
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}

type CreateLikeNotificationEvent struct {
	// Define fields for like notification event
	PostID  uint `json:"post_id"`
	UserID  uint `json:"user_id"`
	LikerID uint `json:"liker_id"`
}

type CreateCommentNotificationEvent struct {
	// Define fields for comment notification event
	PostID      uint   `json:"post_id"`
	UserID      uint   `json:"user_id"`
	CommenterID uint   `json:"commenter_id"`
	Comment     string `json:"comment"`
}
