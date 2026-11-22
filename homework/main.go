package main

import (
	"log"

	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/config"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/controllers"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/routes"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载环境变量
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Warn("Failed to load config file: %v", err)
	}

	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting golang_gin_gorm_blog_system application...")

	// 初始化数据库
	db, err := config.GetConfigDB()
	if err != nil || db == nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	// 初始化全局 service，避免控制器中直接使用 nil 指针。
	controllers.UserService = services.NewUserService(db)
	controllers.PostService = services.NewPostService(db)
	controllers.CommentService = services.NewCommentService(db)

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	logrus.WithField("port", port).Info("Server starting")

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
