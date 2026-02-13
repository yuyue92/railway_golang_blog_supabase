package controllers

import (
	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"

	"github.com/gin-gonic/gin"
)

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// GetCommentsByPostID 获取文章的所有评论
func GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("post_id")

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// 查询评论列表，预加载用户信息
	var comments []models.Comment
	if err := config.DB.Where("post_id = ?", postID).
		Preload("User").
		Order("created_at desc").
		Find(&comments).Error; err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	utils.Success(c, comments)
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID := c.Param("post_id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "评论创建失败")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	utils.SuccessWithMessage(c, "评论创建成功", comment)
}
