package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/services"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/utils"
)

var CommentService *services.CommentService

type CommentController struct{}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	_, err = UserService.GetUserByID(userID.(uint))
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}
	// 检查文章是否存在
	_, err = PostService.GetPostByID(postID)
	if err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	ct, err := CommentService.CreateComment(comment)
	if err != nil {
		utils.InternalServerError(c, "Failed to create comment")
		return
	}

	utils.Success(c, ct)
}

// GetComments 获取文章的评论列表
func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	_, err = PostService.GetPostByID(postID)
	if err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	var comments []models.Comment

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	comments, total, err := CommentService.GetComments(postID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, "Failed to get comments")
		return
	}
	utils.Success(c, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteComment 删除评论
// func (cc *CommentController) DeleteComment(c *gin.Context) {

// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		utils.Unauthorized(c, "User not authenticated")
// 		return
// 	}
// 	_, err := UserService.GetUserByID(userID.(uint))
// 	if err != nil {
// 		utils.NotFound(c, "User not found")
// 		return
// 	}
// 	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		utils.BadRequest(c, "Invalid comment ID")
// 		return
// 	}
// 	comment, err := CommentService.GetCommentByID(commentID)
// 	if err != nil {
// 		utils.NotFound(c, "Comment not found")
// 		return
// 	}
// 	if comment.UserID != userID.(uint) {
// 		utils.Unauthorized(c, "User not authorized")
// 		return
// 	}
// 	err = CommentService.DeleteComment(commentID)
// 	if err != nil {
// 		utils.InternalServerError(c, "Failed to delete comment")
// 		return
// 	}
// 	utils.Success(c, "Comment deleted")
// }
