package models

import (
	"time"

	"gorm.io/gorm"
)
// 定义Post模型
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

// 定义CreatePostRequest模型
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=1"`
}

// 定义UpdatePostRequest模型
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=100"`
	Content string `json:"content" binding:"omitempty,min=1"`
}

