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

type Post struct {
	gorm.Model
	UserID   uint      `gorm:"not null"`
	Content  string    `gorm:"type:text;not null"`
	Comments []Comment `gorm:"foreignKey:PostID"`
	Likes    []Like    `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	PostID  uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
}

type Like struct {
	gorm.Model
	PostID uint `gorm:"not null"`
	UserID uint `gorm:"not null"`
}
