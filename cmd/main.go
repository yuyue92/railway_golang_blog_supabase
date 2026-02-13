package main

import (
	"golang_blog/config"
	"golang_blog/routes"
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用系统环境变量")
	}
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()
	r.Use(gin.Recovery()) // ✅ 防止 panic 变 502

	// 配置CORS - 关键是要包含Authorization
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源，生产环境改为具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,     // 当AllowOrigins为*时，必须为false
		MaxAge:           12 * 3600, // 预检请求缓存时间
	}
	r.Use(cors.New(config))

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
