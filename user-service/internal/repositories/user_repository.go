package repositories

import (
	"github.com/user-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindAll() ([]models.User, error)
	UpdateUser(tx *gorm.DB, user *models.User) error
	FollowUser(tx *gorm.DB, follow *models.UserFollower) error
	GetFollowers(userID uint) ([]models.User, error)
	GetFollowing(userID uint) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Profile").Preload("Settings").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Profile").Preload("Settings").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUser(tx *gorm.DB, user *models.User) error {
	return tx.Save(user).Error
}

func (r *userRepository) FollowUser(tx *gorm.DB, follow *models.UserFollower) error {
	return tx.Create(follow).Error
}

func (r *userRepository) GetFollowers(userID uint) ([]models.User, error) {
	var followers []models.User
	if err := r.db.Joins("JOIN user_followers ON user_followers.follower_id = users.id").
		Where("user_followers.followee_id = ?", userID).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}

func (r *userRepository) GetFollowing(userID uint) ([]models.User, error) {
	var followees []models.User
	if err := r.db.Joins("JOIN user_followers ON user_followers.followee_id = users.id").
		Where("user_followers.follower_id = ?", userID).Find(&followees).Error; err != nil {
		return nil, err
	}
	return followees, nil
}
