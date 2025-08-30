// package api 存放与 API 路由和中间件相关的代码。
package api

import (
	"github.com/KeLes-Coding/gopress/internal/api/handler"
	"github.com/KeLes-Coding/gopress/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 函数用于注册项目的所有 API 路由。
func RegisterRoutes(r *gin.Engine) {
	// 创建一个 API v1 版本的路由分组。
	// 在 URL 路径中保留 /v1 是 RESTful API 的标准实践，这对于 API 版本管理很有好处。
	apiV1Group := r.Group("/api/v1")

	// 实例化各个 handler
	userHandler := handler.NewUserHandler() // <--- 修改实例化方式
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()
	postHandler := handler.NewPostHandler()

	// 公共路由组（无需认证）
	{
		// 注册用户注册接口
		// POST /api/v1/signup
		apiV1Group.POST("/signup", userHandler.SignUpHandler)
		// 注册用户登录接口
		// POST /api/v1/login
		apiV1Group.POST("/login", userHandler.LoginHandler)
	}

	// 认证路由组（需要 JWT 认证）
	authGroup := apiV1Group.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware()) // 应用 JWT 认证中间件
	{
		// 注册获取当前用户信息的接口
		authGroup.GET("/me", userHandler.GetMyProfileHandler)

		// 为后台管理接口创建一个专门的路由组 /admin
		adminGroup := authGroup.Group("/admin")
		{
			// 分类 (Category) 相关路由
			categoryGroup := adminGroup.Group("/categories")
			{
				categoryGroup.POST("", categoryHandler.CreateCategoryHandler)       // 创建分类: POST /api/v1/admin/categories
				categoryGroup.GET("", categoryHandler.ListCategoriesHandler)        // 获取分类列表: GET /api/v1/admin/categories
				categoryGroup.PUT("/:id", categoryHandler.UpdateCategoryHandler)    // 更新分类: PUT /api/v1/admin/categories/:id
				categoryGroup.DELETE("/:id", categoryHandler.DeleteCategoryHandler) // 删除分类: DELETE /api/v1/admin/categories/:id
			}

			// 标签 (Tag) 相关路由
			tagGroup := adminGroup.Group("/tags")
			{
				tagGroup.POST("", tagHandler.CreateTagHandler)       // 创建标签: POST /api/v1/admin/tags
				tagGroup.GET("", tagHandler.ListTagsHandler)         // 获取标签列表: GET /api/v1/admin/tags
				tagGroup.PUT("/:id", tagHandler.UpdateTagHandler)    // 更新标签: PUT /api/v1/admin/tags/:id
				tagGroup.DELETE("/:id", tagHandler.DeleteTagHandler) // 删除标签: DELETE /api/v1/admin/tags/:id
			}

			// 文章 (Post) 相关路由
			postGroup := adminGroup.Group("/posts")
			{
				postGroup.POST("", postHandler.CreatePostHandler)       // 创建文章: POST /api/v1/admin/posts
				postGroup.GET("", postHandler.ListPostsHandler)         // 获取文章列表: GET /api/v1/admin/posts
				postGroup.GET("/:id", postHandler.GetPostHandler)       // 获取单篇文章: GET /api/v1/admin/posts/:id
				postGroup.PUT("/:id", postHandler.UpdatePostHandler)    // 更新文章: PUT /api/v1/admin/posts/:id
				postGroup.DELETE("/:id", postHandler.DeletePostHandler) // 删除文章: DELETE /api/v1/admin/posts/:id
			}
		}
	}
}
