package routes

import (
	"context"
	"golang_blog/config"
	"golang_blog/controllers"
	"golang_blog/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 使用日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})
	r.GET("/ready", func(c *gin.Context) {
		db := config.GetDB()
		if db == nil {
			c.JSON(503, gin.H{"status": "not_ready", "reason": "db is nil"})
			return
		}
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(503, gin.H{"status": "not_ready", "reason": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(503, gin.H{"status": "not_ready", "reason": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由（公开）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 需要认证的路由
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authenticated.GET("/profile", controllers.GetProfile)

			// 文章管理（需要认证）
			authenticated.POST("/posts", controllers.CreatePost)
			authenticated.PUT("/posts/:id", controllers.UpdatePost)
			authenticated.DELETE("/posts/:id", controllers.DeletePost)

			// 评论管理（需要认证）
			authenticated.POST("/posts/:post_id/comments", controllers.CreateComment)
		}

		// 文章路由（公开）
		posts := v1.Group("/posts")
		{
			posts.GET("", controllers.GetPosts)
			posts.GET("/:id", controllers.GetPostByID)
		}

		// 评论路由（公开）
		comments := v1.Group("/comments")
		{
			comments.GET("/post/:post_id", controllers.GetCommentsByPostID)
		}
	}
}
