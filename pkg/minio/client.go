package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"our_records/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

// Init 初始化 MinIO 客户端
func Init() error {
	cfg := config.AppConfig.Minio

	var err error
	Client, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("创建 MinIO 客户端失败：%w", err)
	}

	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return fmt.Errorf("检查存储桶失败：%w", err)
	}

	if !exists {
		err = Client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建存储桶失败：%w", err)
		}
		log.Printf("MinIO 存储桶 %s 创建成功", cfg.Bucket)
	} else {
		log.Printf("MinIO 存储桶 %s 已存在", cfg.Bucket)
	}

	log.Println("MinIO 客户端初始化成功")
	return nil
}

// UploadFile 上传文件到 MinIO，返回相对路径
func UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	cfg := config.AppConfig.Minio

	_, err := Client.PutObject(ctx, cfg.Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// 只返回相对路径，不返回完整 URL
	return objectName, nil
}

// GetPresignedURL 获取签名 URL
func GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	cfg := config.AppConfig.Minio

	// 生成 GET 请求的签名 URL
	url, err := Client.PresignedGetObject(ctx, cfg.Bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

// DeleteFile 删除文件
func DeleteFile(ctx context.Context, objectName string) error {
	return Client.RemoveObject(ctx, config.AppConfig.Minio.Bucket, objectName, minio.RemoveObjectOptions{})
}

// FileExists 检查文件是否存在
func FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := Client.StatObject(ctx, config.AppConfig.Minio.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
