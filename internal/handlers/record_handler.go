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

// GetRecords 获取所有记录（带签名 URL，支持可见性过滤）
func GetRecords(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	// 获取用户所在的群组 ID 列表
	userGroupIDs := getUserGroupIDs(userID)

	var records []models.Record
	query := models.DB.Preload("User")

	// 可见性过滤
	visibility := c.Query("visibility")
	if visibility == "private" {
		// 仅查看公开的和用户所在群组的私有记录
		if len(userGroupIDs) > 0 {
			query = query.Where("visibility = ? OR (visibility = ? AND group_id IN ?)",
				models.VisibilityPublic, models.VisibilityPrivate, userGroupIDs)
		} else {
			// 用户没有加入任何群组，只能看公开的
			query = query.Where("visibility = ?", models.VisibilityPublic)
		}
	} else {
		// 默认查看所有有权限的记录
		if len(userGroupIDs) > 0 {
			query = query.Where("visibility = ? OR user_id = ? OR group_id IN ?",
				models.VisibilityPublic, userID, userGroupIDs)
		} else {
			// 用户没有加入任何群组，只能看公开的和自己的
			query = query.Where("visibility = ? OR user_id = ?",
				models.VisibilityPublic, userID)
		}
	}

	if recordType := c.Query("type"); recordType != "" {
		query = query.Where("type = ?", recordType)
	}
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}
	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
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
			"visibility": record.Visibility,
			"group_id":   record.GroupID,
			"user_id":    record.UserID,
			"username":   "",
			"title":      record.Title,
			"content":    record.Content,
			"tags":       record.Tags,
			"created_at": record.CreatedAt,
			"updated_at": record.UpdatedAt,
		}

		if record.User != nil {
			recordData["username"] = record.User.Username
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

		recordsWithURLs = append(recordsWithURLs, recordData)
	}

	c.JSON(http.StatusOK, recordsWithURLs)
}

// getUserGroupIDs 获取用户所在的所有群组 ID
func getUserGroupIDs(userID uint) []uint {
	var userGroups []models.UserGroup
	models.DB.Where("user_id = ?", userID).Pluck("group_id", &userGroups)
	ids := make([]uint, 0, len(userGroups))
	for _, ug := range userGroups {
		ids = append(ids, ug.GroupID)
	}
	return ids
}

// GetRecord 获取单条记录（带签名 URL）
func GetRecord(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

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

	// 权限验证：检查用户是否有权限查看此记录
	if record.Visibility == models.VisibilityPrivate && (record.UserID == nil || *record.UserID != userID) {
		// 私有记录，检查用户是否在群组中
		userGroupIDs := getUserGroupIDs(userID)
		hasPermission := false
		if record.GroupID != nil {
			for _, gid := range userGroupIDs {
				if gid == *record.GroupID {
					hasPermission = true
					break
				}
			}
		}
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此记录"})
			return
		}
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

// CreateRecord 创建记录
func CreateRecord(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	var record models.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置用户 ID
	record.UserID = &userID

	// 默认为公开
	if record.Visibility == "" {
		record.Visibility = models.VisibilityPublic
	}

	if err := models.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 预加载用户信息
	models.DB.Preload("User").First(&record, record.ID)

	c.JSON(http.StatusCreated, record)
}

// UpdateRecord 更新记录
func UpdateRecord(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

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

	// 只有记录所有者才能更新
	if record.UserID == nil || *record.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此记录"})
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
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

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

	// 只有记录所有者才能删除
	if record.UserID == nil || *record.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此记录"})
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
