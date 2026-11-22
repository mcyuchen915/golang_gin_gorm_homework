package services

import (
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// 创建评论
func (ps *CommentService) CreateComment(comment models.Comment) (*models.Comment, error) {

	if err := ps.db.Create(&comment).Error; err != nil {

		return nil, err
	}

	// 预加载用户信息
	ps.db.Preload("User").First(&comment, comment.ID)

	return &comment, nil
}

// 获取评论
func (ps *CommentService) GetCommentByID(id uint64) (*models.Comment, error) {
	var comment models.Comment
	if err := ps.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetComments 获取文章的评论列表
func (s *CommentService) GetComments(postID uint64, page, pageSize int) ([]models.Comment, uint64, error) {
	var comments []models.Comment

	// 获取总数
	var total int64
	s.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)
	if total > 0 {
		offset := (page - 1) * pageSize
		// 查询评论列表，预加载用户信息
		if err := s.db.Preload("User").
			Where("post_id = ?", postID).
			Order("created_at ASC").
			Limit(pageSize).
			Offset(offset).
			Find(&comments).Error; err != nil {
			return nil, 0, err
		}
	}

	return comments, uint64(total), nil
}

// 删除评论
func (s *CommentService) DeleteComment(id uint64) error {
	if err := s.db.Delete(&models.Comment{}, id).Error; err != nil {
		return err
	}
	return nil
}
