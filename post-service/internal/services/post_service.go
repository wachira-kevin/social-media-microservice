package services

import (
	"post-service/internal/models"
	"post-service/internal/repositories"
)

type PostService struct {
	PostRepository repositories.PostRepository
}

func NewPostService(postRepository repositories.PostRepository) *PostService {
	return &PostService{
		PostRepository: postRepository,
	}
}

func (s *PostService) CreatePost(post *models.Post) error {
	return s.PostRepository.CreatePost(post)
}

func (s *PostService) UpdatePost(post *models.Post) error {
	return s.PostRepository.UpdatePost(post)
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	return s.PostRepository.GetPostByID(id)
}

func (s *PostService) GetPostsByUserID(userID uint) ([]models.Post, error) {
	return s.PostRepository.GetPostsByUserID(userID)
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.PostRepository.GetPosts()
}

func (s *PostService) LikePost(postID, userID uint) error {
	like := models.Like{PostID: postID, UserID: userID}
	return s.PostRepository.LikePost(&like)
}

func (s *PostService) CommentOnPost(comment *models.Comment) error {
	return s.PostRepository.CommentOnPost(comment)
}
