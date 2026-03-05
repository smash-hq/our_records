package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
)

// GetRecords 获取所有记录（带签名 URL）
func GetRecords(c *gin.Context) {
	var records []models.Record
	query := models.DB

	if recordType := c.Query("type"); recordType != "" {
		query = query.Where("type = ?", recordType)
	}
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	if err := query.Order("created_at DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recordsWithURLs := make([]gin.H, 0, len(records))
	for _, record := range records {
		recordData := gin.H{
			"id":         record.ID,
			"type":       record.Type,
			"title":      record.Title,
			"content":    record.Content,
			"tags":       record.Tags,
			"created_at": record.CreatedAt,
			"updated_at": record.UpdatedAt,
		}

		if record.MediaPath != "" && (record.Type == "image" || record.Type == "audio" || record.Type == "video") {
			paths := record.GetMediaPaths()
			signedURLs := make([]string, 0, len(paths))
			for _, p := range paths {
				signedURL, err := minioClient.GetPresignedURL(context.Background(), p, 24*time.Hour)
				if err == nil {
					signedURLs = append(signedURLs, signedURL)
				}
			}
			recordData["media_paths"] = signedURLs
			recordData["media_path"] = signedURLs[0] // 兼容旧字段
		} else {
			recordData["media_paths"] = []string{}
			recordData["media_path"] = ""
		}

		recordsWithURLs = append(recordsWithURLs, recordData)
	}

	c.JSON(http.StatusOK, recordsWithURLs)
}

// GetRecord 获取单条记录（带签名 URL）
func GetRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	recordData := gin.H{
		"id":         record.ID,
		"type":       record.Type,
		"title":      record.Title,
		"content":    record.Content,
		"tags":       record.Tags,
		"created_at": record.CreatedAt,
		"updated_at": record.UpdatedAt,
	}

	if record.MediaPath != "" && (record.Type == "image" || record.Type == "audio" || record.Type == "video") {
		paths := record.GetMediaPaths()
		signedURLs := make([]string, 0, len(paths))
		for _, p := range paths {
			signedURL, err := minioClient.GetPresignedURL(context.Background(), p, 24*time.Hour)
			if err == nil {
				signedURLs = append(signedURLs, signedURL)
			}
		}
		recordData["media_paths"] = signedURLs
		recordData["media_path"] = signedURLs[0]
	} else {
		recordData["media_paths"] = []string{}
		recordData["media_path"] = ""
	}

	c.JSON(http.StatusOK, recordData)
}

// CreateRecord 创建记录（纯文字）
func CreateRecord(c *gin.Context) {
	var record models.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// UpdateRecord 更新记录
func UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.ID = uint(id)
	if err := models.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord 删除记录（逻辑删除）
func DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if err := models.DB.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetMediaURL 获取媒体签名 URL
func GetMediaURL(c *gin.Context) {
	objectPath := c.Query("path")
	if objectPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少路径参数"})
		return
	}

	expiry := 24 * time.Hour
	if exp := c.Query("expiry"); exp != "" {
		if h, err := strconv.Atoi(exp); err == nil && h > 0 {
			expiry = time.Duration(h) * time.Hour
		}
	}

	signedURL, err := minioClient.GetPresignedURL(context.Background(), objectPath, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成签名 URL 失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": signedURL, "expiry": expiry})
}

// parseMediaPaths 解析媒体路径
func parseMediaPaths(mediaPath string) []string {
	if mediaPath == "" {
		return []string{}
	}
	paths := strings.Split(mediaPath, ",")
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
