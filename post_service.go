package services

import (
	"errors"

	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (ps *PostService) CreatePost(post models.Post) (*models.Post, error) {

	if err := ps.db.Create(&post).Error; err != nil {

		return nil, err
	}

	// 预加载用户信息
	ps.db.Preload("User").First(&post, post.ID)

	return &post, nil
}

// GetPostByID 获取文章
func (s *PostService) GetPostByID(postID uint64) (*models.Post, error) {
	var post models.Post

	if err := s.db.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(id uint64, req models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	// 更新文章
	post.Title = req.Title
	post.Content = req.Content
	if err := s.db.Save(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

// 文章分页
func (s *PostService) GetPosts(page, pageSize int) ([]models.Post, uint64, error) {
	if s == nil || s.db == nil {
		return nil, 0, errors.New("post service is not initialized")
	}

	var posts []models.Post

	// 获取总数
	var total int64
	if err := s.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total > 0 {
		offset := (page - 1) * pageSize
		if err := s.db.Preload("User").
			Order("created_at DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&posts).Error; err != nil {
			return nil, 0, err
		}
	}

	return posts, uint64(total), nil
}

// 删除文章
func (s *PostService) DeletePost(id uint64) error {
	if err := s.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}
	return nil
}
