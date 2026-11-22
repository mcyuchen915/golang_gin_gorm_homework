package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/controllers"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Recovery())

	// 创建控制器实例
	userController := &controllers.UserController{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			// 用户信息
			authenticated.GET("/profile", userController.GetUser)
			// 文章相关路由
			posts := authenticated.Group("/posts")
			{
				posts.POST("", postController.CreatePost)
				posts.PUT("/:id", postController.UpdatePost)
				posts.DELETE("/:id", postController.DeletePost)
			}

			// 评论相关路由
			comments := authenticated.Group("/posts/:post_id/comments")
			{
				comments.POST("", commentController.CreateComment)
			}
			comments = authenticated.Group("/comment/:post_id/comments/:comment_id")
			// comments.DELETE("", commentController.DeleteComment)
		}

		// 公开路由（无需认证）
		public := api.Group("")
		{
			// 文章公开路由
			public.GET("/posts", postController.GetPosts)
			public.GET("/posts/:id", postController.GetPost)
		}

		// 评论公开路由（单独分组避免路由冲突）
		comments := api.Group("/comments")
		{
			comments.GET("/post/:post_id", commentController.GetComments)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "gin_gorm_Blog API is running",
		})
	})

	return r
}
