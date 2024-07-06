package services

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/user-service/config"
	"github.com/user-service/internal/models"
	"github.com/user-service/internal/publishers"
	"github.com/user-service/internal/repositories"
	"github.com/user-service/pkg/utils"
	"gorm.io/gorm"
	"log"
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
	userRepository  repositories.UserRepository
	keycloakService *KeycloakService
	conn            *amqp.Connection
	db              *gorm.DB
}

func (u userService) FollowUser(followerID uint, followeeID uint) error {
	// get follower and followee
	follower, err := u.userRepository.FindByID(followerID)
	if err != nil {
		return err
	}
	followee, err := u.userRepository.FindByID(followeeID)
	if err != nil {
		return err
	}

	// Check if the user already follows the followee
	exists, err := u.userRepository.DoesFollowExist(followerID, followeeID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user %d already follows user %d", followerID, followeeID)
	}

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

	// adding follow notification to rabbit to be consumed by notification
	if followee.Settings.PushNotifications {
		notification := &models.NotificationSchema{
			UserID:           followeeID,
			NotificationType: "push",
			MessageType:      "New Follower",
			Message:          fmt.Sprintf("You have a new follower: %s", follower.Username),
			Email:            followee.Email,
		}
		err := publishers.PublishFollowingMessage(notification, u.conn)
		if err != nil {
			log.Fatalf("Error publishing new follower message: %v", err)
		}
	} else {
		notification := &models.NotificationSchema{
			UserID:           followeeID,
			NotificationType: "email",
			MessageType:      "follow",
			Message:          fmt.Sprintf("You have a new follower: %s", follower.Username),
			Email:            followee.Email,
		}
		err := publishers.PublishFollowingMessage(notification, u.conn)
		if err != nil {
			log.Fatalf("Error publishing new follower message: %v", err)
		}
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
	// check if user exists
	userExists, err := u.userRepository.UserExists(userInput.Username, userInput.Email)
	if err != nil {
		return nil, err
	}
	if userExists {
		return nil, fmt.Errorf("conflict")
	}

	// make request to keycloak
	keycloakID, err := u.keycloakService.CreateUserInKeycloak(userInput)
	if err != nil {
		return nil, err
	}

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
		KeycloakId: keycloakID,
		Username:   userInput.Username,
		Email:      userInput.Email,
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

func NewUserService(db *gorm.DB, cfg *config.Config, conn *amqp.Connection) UserService {
	return &userService{
		userRepository:  repositories.NewUserRepository(db),
		keycloakService: NewKeycloakService(cfg),
		conn:            conn,
		db:              db,
	}
}
