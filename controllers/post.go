package controllers

import (
	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePostRequest 创建文章请求结构
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var posts []models.Post
	var total int64

	// 统计总数
	config.DB.Model(&models.Post{}).Count(&total)

	// 查询文章列表，预加载用户信息
	if err := config.DB.Preload("User").
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	utils.Success(c, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetPostByID 获取文章详情
func GetPostByID(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	utils.Success(c, post)
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
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

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		utils.InternalServerError(c, "文章创建失败")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)

	utils.SuccessWithMessage(c, "文章创建成功", post)
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var req CreatePostRequest
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

	// 查找文章
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "无权限修改此文章")
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content

	if err := config.DB.Save(&post).Error; err != nil {
		utils.InternalServerError(c, "文章更新失败")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)

	utils.SuccessWithMessage(c, "文章更新成功", post)
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	// 查找文章
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "无权限删除此文章")
		return
	}

	// 删除文章（软删除）
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "文章删除失败")
		return
	}

	utils.SuccessWithMessage(c, "文章删除成功", nil)
}
