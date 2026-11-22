package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/services"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/utils"
)

type PostController struct{}

var PostService *services.PostService

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}
	crepost, err := PostService.CreatePost(post)
	if err != nil {
		utils.InternalServerError(c, "Failed to create post")
		return
	}
	utils.Success(c, crepost)
}

// GetPosts 获取文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询文章列表，预加载用户信息
	posts, total, err := PostService.GetPosts(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, "Failed to get posts")
		return
	}

	utils.Success(c, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPost 获取单个文章详情
func (pc *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}
	post, err := PostService.GetPostByID(postID)
	if err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	utils.Success(c, post)
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	post, err := PostService.GetPostByID(postID)
	if err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only update your own posts")
		return
	}

	post, err = PostService.UpdatePost(postID, req)
	if err != nil {
		utils.InternalServerError(c, "Failed to update post")
		return
	}
	utils.Success(c, post)
}

// DeletePost 删除文章（软删除）
func (pc *PostController) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	post, err := PostService.GetPostByID(postID)
	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only delete your own posts")
		return
	}

	// 删除文章（软删除）
	if err := PostService.DeletePost(postID); err != nil {
		utils.InternalServerError(c, "Failed to delete post")
		return
	}

	utils.Success(c, gin.H{"message": "Post deleted successfully"})
}
