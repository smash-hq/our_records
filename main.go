package main

import (
	"log"

	"our_records/internal/config"
	"our_records/internal/middleware"
	"our_records/internal/models"
	"our_records/internal/routes"
	"our_records/pkg/minio"
)

func main() {
	// 初始化配置
	if err := config.Init("application.yaml"); err != nil {
		log.Fatalf("配置加载失败：%v", err)
	}

	// 初始化 JWT
	middleware.InitJWT()

	// 初始化数据库
	db := config.AppConfig.Database
	if err := models.InitDB(db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode); err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
	}

	// 初始化 MinIO 客户端
	if err := minio.Init(); err != nil {
		log.Fatalf("MinIO 初始化失败：%v", err)
	}

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
