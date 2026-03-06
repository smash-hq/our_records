package routes

import (
	"net/http"
	"our_records/internal/handlers"
	"our_records/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 认证相关路由（公开）
	api := r.Group("/api")
	{
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关
		auth.GET("/user", handlers.GetCurrentUser)
		auth.PUT("/user/password", handlers.ChangePassword)
		auth.POST("/user/avatar", handlers.UploadAvatar)

		// 群组相关路由
		auth.POST("/groups", handlers.CreateGroup)
		auth.GET("/groups", handlers.GetGroups)
		auth.GET("/groups/:id", handlers.GetGroup)
		auth.PUT("/groups/:id", handlers.UpdateGroup)
		auth.DELETE("/groups/:id", handlers.DeleteGroup)
		auth.POST("/groups/:id/avatar", handlers.UploadGroupAvatar)
		auth.POST("/groups/:id/members", handlers.AddGroupMember)
		auth.DELETE("/groups/:id/members", handlers.RemoveGroupMember)
		auth.GET("/groups/:id/members", handlers.GetGroupMembers)
		auth.POST("/groups/:id/leave", handlers.LeaveGroup)
		auth.GET("/users/search", handlers.SearchUsers)

		// 记录相关路由
		records := auth.Group("/records")
		{
			records.POST("", handlers.CreateRecord)
			records.GET("", handlers.GetRecords)
			records.GET("/:id", handlers.GetRecord)
			records.PUT("/:id", handlers.UpdateRecord)
			records.DELETE("/:id", handlers.DeleteRecord)
		}

		// 文件上传路由
		auth.POST("/upload", handlers.UploadFile)
	}

	return r
}
