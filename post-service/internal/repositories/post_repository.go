package repositories

import (
	"gorm.io/gorm"
	"post-service/internal/models"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	UpdatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	GetPostsByUserID(userID uint) ([]models.Post, error)
	GetPosts() ([]models.Post, error)
	LikePost(like *models.Like) error
	CommentOnPost(comment *models.Comment) error
}

type postRepository struct {
	db *gorm.DB
}

func (r *postRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Preload("Comments").Preload("Likes").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) UpdatePost(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Comments").Preload("Likes").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostsByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Preload("Comments").Preload("Likes").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) LikePost(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *postRepository) CommentOnPost(comment *models.Comment) error {
	return r.db.Create(comment).Error
}
