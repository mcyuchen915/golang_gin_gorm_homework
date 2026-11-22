一个使用Go语言、Gin框架和GORM库开发的个人博客系统后端API，使用类似Java MVC的架构，使用JWT进行身份认证和授权，使用Logrus进行日志记录，使用SQLite作为数据库。

一、功能特性
1、用户注册和登录
2、JWT认证和授权
3、日志记录
4、文章的CRUD操作
5、评论功能
6、统一的错误处理
7、数据库自动迁移
二、技术栈
1、Go 1.27.3
2、Gin - Web框架
3、GORM - ORM库
4、SQLite - 数据库
5、viper - 配置管理
6、JWT - 身份认证
7、Logrus - 日志库
8、bcrypt - 密码加密
三、依赖安装步骤和启动方式
1、go mod tidy
2、go run main.go 运行main.go文件，进行启动

说明：
1、路由配置在router.go文件中
2、server、database, jwt等配置在config.yaml中
3、服务配置数据库获取和初始化配置等在config.go文件中


测试用例
1、测试用例在test目录下
router路由路径如下：

方法	路径	Handler	认证
POST	/api/v1/auth/register	                userController.Register	        无需认证
POST	/api/v1/auth/login	                    userController.Login	        无需认证
GET	    /api/v1/profile	                        userController.GetUser	        需要认证
POST	/api/v1/posts	                        postController.CreatePost	    需要认证
PUT	    /api/v1/posts/:id	                    postController.UpdatePost	    需要认证
DELETE	/api/v1/posts/:id	                    postController.DeletePost	    需要认证
POST	/api/v1/posts/:post_id/comments	        commentController.CreateComment	需要认证
DELETE	/comment/:post_id/comments/:comment_id	commentController.DeleteComment	需要认证
GET	    /api/v1/posts	                        postController.GetPosts	        无需认证
GET 	/api/v1/posts/:id	                    postController.GetPost	        无需认证
GET	    /api/v1/comments/post/:post_id	        commentController.GetComments	无需认证
GET	    /health	匿名函数	无需认证

测试用例：
使用postman进行测试
1、注册请求：http://localhost:8080/api/v1/auth/register

POST /user
Content-Type: application/json

Request Body:
{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "password": "123456"
}
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDM2MywibmJmIjoxNzk1MzQ3OTYzLCJpYXQiOjE3OTUzNDc5NjN9.nK342Kk-t4SGOVdi0Kqdv9TI8bCWkPRLKd34GVjX7k4",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        }
    }
}

2、登录请求：http://localhost:8080/api/v1/auth/login

POST /user/login
Content-Type: application/json

Request Body:
{
    "username": "zhangsan",
    "password": "123456"
}
Reponse body:
{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        }
    }
}

3、获取用户信息：http://localhost:8080/api/v1/profile

GET /user/profile
Request Header:
Authorization: Bearer <eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs>

Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 19,
        "username": "zhangsan",
        "email": "zhangsan@example.com",
        "created_at": "2026-11-22T19:46:03.556214+08:00",
        "updated_at": "2026-11-22T19:46:03.556214+08:00"
    }
}

4、创建文章：http://localhost:8080/api/v1/posts

POST /posts
Content-Type: application/json
Request Header:
Authorization: Bearer <eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs>

Request Body:
{
    "title": "My First Post",
    "content": "This is my first post."
}
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "title": "My First Post",
        "content": "This is my first post.",
        "user_id": 19,
        "created_at": "2026-11-22T21:13:52.8190035+08:00",
        "updated_at": "2026-11-22T21:13:52.8190035+08:00",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        }
    }
}

5、获取文章列表：http://localhost:8080/api/v1/posts

GET
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "page": 1,
        "page_size": 10,
        "posts": [
            {
                "id": 2,
                "title": "My First Post",
                "content": "This is my first post.",
                "user_id": 19,
                "created_at": "2026-11-22T21:16:12.5650549+08:00",
                "updated_at": "2026-11-22T21:16:12.5650549+08:00",
                "user": {
                    "id": 19,
                    "username": "zhangsan",
                    "email": "zhangsan@example.com",
                    "created_at": "2026-11-22T19:46:03.556214+08:00",
                    "updated_at": "2026-11-22T19:46:03.556214+08:00"
                }
            },
            {
                "id": 1,
                "title": "My First Post",
                "content": "This is my first post.",
                "user_id": 19,
                "created_at": "2026-11-22T21:13:52.8190035+08:00",
                "updated_at": "2026-11-22T21:13:52.8190035+08:00",
                "user": {
                    "id": 19,
                    "username": "zhangsan",
                    "email": "zhangsan@example.com",
                    "created_at": "2026-11-22T19:46:03.556214+08:00",
                    "updated_at": "2026-11-22T19:46:03.556214+08:00"
                }
            }
        ],
        "total": 2
    }
}

6、 获取文章详情：http://localhost:8080/api/v1/posts/1

GET /posts/1
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "title": "My First Post",
        "content": "This is my first post.",
        "user_id": 19,
        "created_at": "2026-11-22T21:13:52.8190035+08:00",
        "updated_at": "2026-11-22T21:13:52.8190035+08:00",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        }
    }
}

7、 更新文章：http://localhost:8080/api/v1/posts/2

PUT /posts/2
Content-Type: application/json
Request Header:
Authorization: Bearer <eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs>

Request Body:
{
    "title": "My second Post",
    "content": "This is my second post."
}
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 2,
        "title": "My second Post",
        "content": "This is my second post.",
        "user_id": 19,
        "created_at": "2026-11-22T21:16:12.5650549+08:00",
        "updated_at": "2026-11-22T21:25:10.8536305+08:00",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        }
    }
}

8、 删除文章：http://localhost:8080/api/v1/posts/3

DELETE /posts/3
Request Header:
Authorization: Bearer <eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs>

Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "message": "Post deleted successfully"
    }
}

9、 创建评论：http://localhost:8080/api/v1/posts/1/comments

POST /posts/1/comments
Content-Type: application/json
Request Header:
Authorization: Bearer <eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwidXNlcm5hbWUiOiJ6aGFuZ3NhbiIsImV4cCI6MTc5NTQzNDc0OCwibmJmIjoxNzk1MzQ4MzQ4LCJpYXQiOjE3OTUzNDgzNDh9.jDmBbgflMHMw03nlQo5v9mgem7Hgppyn2mZrPtihXGs>

Request Body:
{
    "content": "This is my first comment."
}
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "content": "This is my first comment.",
        "user_id": 19,
        "post_id": 1,
        "created_at": "2026-11-22T21:31:36.9193367+08:00",
        "updated_at": "2026-11-22T21:31:36.9193367+08:00",
        "user": {
            "id": 19,
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "created_at": "2026-11-22T19:46:03.556214+08:00",
            "updated_at": "2026-11-22T19:46:03.556214+08:00"
        },
        "post": {
            "id": 0,
            "title": "",
            "content": "",
            "user_id": 0,
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "user": {
                "id": 0,
                "username": "",
                "email": "",
                "created_at": "0001-01-01T00:00:00Z",
                "updated_at": "0001-01-01T00:00:00Z"
            }
        }
    }
}

10、获取评论列表：http://localhost:8080/api/v1/comments/post/1
GET /comments/post/1
Response Body:
{
    "code": 200,
    "message": "success",
    "data": {
        "comments": [
            {
                "id": 1,
                "content": "This is my first comment.",
                "user_id": 19,
                "post_id": 1,
                "created_at": "2026-11-22T21:31:36.9193367+08:00",
                "updated_at": "2026-11-22T21:31:36.9193367+08:00",
                "user": {
                    "id": 19,
                    "username": "zhangsan",
                    "email": "zhangsan@example.com",
                    "created_at": "2026-11-22T19:46:03.556214+08:00",
                    "updated_at": "2026-11-22T19:46:03.556214+08:00"
                },
                "post": {
                    "id": 0,
                    "title": "",
                    "content": "",
                    "user_id": 0,
                    "created_at": "0001-01-01T00:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "user": {
                        "id": 0,
                        "username": "",
                        "email": "",
                        "created_at": "0001-01-01T00:00:00Z",
                        "updated_at": "0001-01-01T00:00:00Z"
                    }
                }
            }
        ],
        "page": 1,
        "page_size": 20,
        "total": 1
    }
}
