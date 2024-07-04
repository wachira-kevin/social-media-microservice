package services

import (
	"github.com/user-service/internal/models"
	"github.com/user-service/internal/repositories"
	"github.com/user-service/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type UserService interface {
	CreateUser(userInput *models.UserCreationSchema) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(userID uint, userInput *models.UserUpdateSchema) (*models.User, error)
	FollowUser(followerID uint, followeeID uint) error
	GetFollowers(userID uint) ([]models.User, error)
	GetFollowing(userID uint) ([]models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
	db             *gorm.DB
}

func (u userService) FollowUser(followerID uint, followeeID uint) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	follow := &models.UserFollower{
		FollowerID: followerID,
		FolloweeID: followeeID,
		FollowedAt: time.Now(),
	}

	if err := u.userRepository.FollowUser(tx, follow); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u userService) GetFollowers(userID uint) ([]models.User, error) {
	return u.userRepository.GetFollowers(userID)
}

func (u userService) GetFollowing(userID uint) ([]models.User, error) {
	return u.userRepository.GetFollowing(userID)
}

func (u userService) CreateUser(userInput *models.UserCreationSchema) (*models.User, error) {
	// parsing date
	dob, err := utils.ParseDate(userInput.DateOfBirth)
	if err != nil {
		return nil, err
	}

	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create User
	user := &models.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		Settings: models.UserSettings{
			EmailNotifications: true,
			PushNotifications:  false,
			SmsNotifications:   false,
		},
	}

	// Create Profile
	profile := &models.Profile{
		FirstName:   userInput.FirstName,
		LastName:    userInput.LastName,
		Bio:         userInput.Bio,
		DateOfBirth: &dob,
		Gender:      userInput.Gender,
	}
	user.Profile = *profile

	// Save User and Profile
	if err := u.userRepository.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) GetUserByID(userID uint) (*models.User, error) {
	return u.userRepository.FindByID(userID)
}

func (u userService) GetAllUsers() ([]models.User, error) {
	return u.userRepository.FindAll()
}

func (u userService) UpdateUser(userID uint, userInput *models.UserUpdateSchema) (*models.User, error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := u.userRepository.FindByID(userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if userInput.Username != "" {
		user.Username = userInput.Username
	}
	if userInput.Email != "" {
		user.Email = userInput.Email
	}
	if userInput.FirstName != "" {
		user.Profile.FirstName = userInput.FirstName
	}
	if userInput.LastName != "" {
		user.Profile.LastName = userInput.LastName
	}
	if userInput.Bio != "" {
		user.Profile.Bio = userInput.Bio
	}
	if userInput.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("02/01/2006", userInput.DateOfBirth)
		if err == nil {
			user.Profile.DateOfBirth = &dateOfBirth
		}
	}
	if userInput.Gender != "" {
		user.Profile.Gender = userInput.Gender
	}
	if userInput.Location != "" {
		user.Profile.Location = userInput.Location
	}
	if userInput.NotificationType != "" {
		switch userInput.NotificationType {
		case "email":
			user.Settings.EmailNotifications = true
			user.Settings.PushNotifications = false
			user.Settings.SmsNotifications = false
		case "push":
			user.Settings.EmailNotifications = false
			user.Settings.PushNotifications = true
			user.Settings.SmsNotifications = false
		case "sms":
			user.Settings.EmailNotifications = false
			user.Settings.PushNotifications = false
			user.Settings.SmsNotifications = true
		}
	}

	if err := u.userRepository.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil

}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		userRepository: repositories.NewUserRepository(db),
		db:             db,
	}
}
