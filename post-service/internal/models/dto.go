package models

import (
	"time"
)

type CreatePost struct {
	UserID  uint   `json:"user_id" binding:"required,numeric"`
	Content string `json:"content" binding:"required,max=250"`
}

type EditPost struct {
	Content string `json:"content" binding:"required,max=250"`
}

type CreateComment struct {
	UserID  uint   `json:"user_id" binding:"required,numeric"`
	Content string `json:"content" binding:"required,max=250"`
}

type User struct {
	ID         uint         `json:"id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  time.Time    `json:"deleted_at"`
	KeycloakId string       `json:"keycloak_id"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	LastLogin  time.Time    `json:"last_login"`
	Profile    profile      `json:"profile"`
	Settings   userSettings `json:"settings"`
}

type profile struct {
	ID                uint      `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
	UserID            uint      `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Bio               string    `json:"bio"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	Gender            string    `json:"gender"`
	Location          string    `json:"location"`
}

type userSettings struct {
	ID                 uint      `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          time.Time `json:"deleted_at"`
	UserID             uint      `json:"user_id"`
	EmailNotifications bool      `json:"email_notifications"`
	PushNotifications  bool      `json:"push_notifications"`
	SmsNotifications   bool      `json:"sms_notifications"`
}
