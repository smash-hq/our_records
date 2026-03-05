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

	// API 路由组
	api := r.Group("/api")
	{
		// 记录相关路由
		records := api.Group("/records")
		{
			records.POST("", handlers.CreateRecord)
			records.GET("", handlers.GetRecords)
			records.GET("/:id", handlers.GetRecord)
			records.PUT("/:id", handlers.UpdateRecord)
			records.DELETE("/:id", handlers.DeleteRecord)
		}

		// 文件上传路由
		api.POST("/upload", handlers.UploadFile)
	}

	return r
}
