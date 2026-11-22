package models

import (
	"time"

	"gorm.io/gorm"
)

// 定义Comment模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
}

// 定义CreateCommentRequest模型
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// 定义CommentResponse模型
type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
