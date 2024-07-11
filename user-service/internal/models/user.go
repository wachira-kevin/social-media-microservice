package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	gorm.Model
	KeycloakId string       `gorm:"unique;not null" json:"keycloak_id,omitempty"`
	Username   string       `gorm:"unique;not null" json:"username"`
	Email      string       `gorm:"unique;not null" json:"email"`
	LastLogin  *time.Time   `json:"last_login,omitempty"`
	Profile    Profile      `gorm:"foreignKey:UserID" json:"profile"`
	Settings   UserSettings `gorm:"foreignKey:UserID" json:"settings"`
}

type Profile struct {
	gorm.Model
	UserID            uint       `gorm:"not null" json:"user_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Bio               string     `json:"bio"`
	ProfilePictureURL string     `json:"profile_picture_url,omitempty"`
	DateOfBirth       *time.Time `json:"date_of_birth"`
	Gender            string     `json:"gender"`
	Location          string     `json:"location"`
}

type UserFollower struct {
	gorm.Model
	FollowerID uint      `gorm:"not null" json:"follower_id"`
	FolloweeID uint      `gorm:"not null" json:"followee_id"`
	FollowedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"followed_at"`
}

type UserSettings struct {
	gorm.Model
	UserID             uint `gorm:"primaryKey" json:"user_id"`
	EmailNotifications bool `gorm:"default:false" json:"email_notifications"`
	PushNotifications  bool `gorm:"default:false" json:"push_notifications"`
	SmsNotifications   bool `gorm:"default:false" json:"sms_notifications"`
}

type UserCreationSchema struct {
	FirstName   string `json:"first_name" binding:"required,alpha"`
	LastName    string `json:"last_name" binding:"required,alpha"`
	Bio         string `json:"bio" binding:"omitempty,max=250"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
	Gender      string `json:"gender" binding:"required,oneof=male female other"`
	Location    string `json:"location" binding:"required"`
	Username    string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password    string `json:"password" binding:"required,min=3,max=20"`
	Email       string `json:"email" binding:"required,email"`
}

type UserUpdateSchema struct {
	FirstName   string `json:"first_name" binding:"omitempty,alpha"`
	LastName    string `json:"last_name" binding:"omitempty,alpha"`
	Bio         string `json:"bio" binding:"omitempty,max=250"`
	DateOfBirth string `json:"date_of_birth" binding:"omitempty"`
	Gender      string `json:"gender" binding:"omitempty,oneof=male female other"`
	Location    string `json:"location" binding:"omitempty"`
	Username    string `json:"username" binding:"omitempty,alphanum,min=3,max=20"`
	Email       string `json:"email" binding:"omitempty,email"`
}

type UserSettingsUpdateSchema struct {
	NotificationType string `json:"notification_type" binding:"omitempty,oneof=sms email push"`
}

type NotificationSchema struct {
	UserID           uint   `json:"user_id"`
	MessageType      string `json:"message_type"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
	Email            string `json:"email"`
}
