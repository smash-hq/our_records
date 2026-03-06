package handlers

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"our_records/internal/config"
	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件到 MinIO，只存储相对路径
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件"})
		return
	}

	// 检查文件大小
	if file.Size > config.AppConfig.Upload.MaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	var recordType models.RecordType
	var objectPrefix string

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		recordType = models.RecordTypeImage
		objectPrefix = "images"
	case ".mp3", ".wav", ".ogg", ".m4a":
		recordType = models.RecordTypeAudio
		objectPrefix = "audios"
	case ".mp4", ".avi", ".mov", ".webm":
		recordType = models.RecordTypeVideo
		objectPrefix = "videos"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 生成新文件名
	newFilename := time.Now().Format("20060102150405") + "_" + file.Filename
	objectName := objectPrefix + "/" + newFilename

	// 打开上传的文件
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	// 读取文件内容
	data := make([]byte, file.Size)
	_, err = openedFile.Read(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确定 Content-Type
	contentType := getContentType(ext)

	// 上传到 MinIO，返回相对路径
	objectPath, err := minioClient.UploadFile(context.Background(), objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO 上传失败：" + err.Error()})
		return
	}

	// 检查是否已有记录（多图上传）
	var record models.Record
	title := c.PostForm("title")
	content := c.PostForm("content")
	tags := c.PostForm("tags")
	visibility := c.PostForm("visibility")
	groupIDStr := c.PostForm("group_id")

	// 尝试查找同标题的记录
	models.DB.Where("title = ? AND type = ?", title, recordType).Order("id DESC").First(&record)

	mediaPath := objectPath
	if record.ID > 0 && record.MediaPath != "" {
		// 追加到现有路径
		mediaPath = record.MediaPath + "," + objectPath
	}

	// 创建或更新记录
	if record.ID > 0 {
		record.MediaPath = mediaPath
		record.Content = content
		record.Tags = tags
		models.DB.Save(&record)
	} else {
		// 解析 group_id
		var groupID *uint
		if groupIDStr != "" {
			var gid uint
			if err := models.DB.Raw("SELECT ? AS group_id", groupIDStr).Scan(&gid).Error; err == nil {
				groupID = &gid
			}
		}

		record = models.Record{
			Type:       recordType,
			Title:      title,
			Content:    content,
			MediaPath:  mediaPath,
			Tags:       tags,
			Visibility: models.Visibility(visibility),
			GroupID:    groupID,
		}
		// 默认为公开
		if record.Visibility == "" {
			record.Visibility = models.VisibilityPublic
		}
		if err := models.DB.Create(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "上传成功",
		"record":      record,
		"storage":     "minio",
		"object_path": objectPath,
	})
}

// getContentType 根据扩展名获取 Content-Type
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".ogg":  "audio/ogg",
		".m4a":  "audio/mp4",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
