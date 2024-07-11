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

type Notification struct {
	gorm.Model
	UserID           uint   `gorm:"not null" json:"user_id"`
	MessageType      string `gorm:"not null" json:"message_type"`
	NotificationType string `gorm:"not null" json:"notification_type"`
	Message          string `gorm:"not null" json:"message"`
	Status           string `gorm:"not null" json:"status"`
}

type NotificationSchema struct {
	UserID           uint   `json:"user_id"`
	MessageType      string `json:"message_type"`
	NotificationType string `json:"notification_type"`
	Email            string `json:"email"`
	Message          string `json:"message"`
}
