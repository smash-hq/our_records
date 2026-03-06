package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
)

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID           uint                    `json:"id"`
	UserID       uint                    `json:"user_id"`
	FromUserID   uint                    `json:"from_user_id"`
	FromUsername string                  `json:"from_username"`
	FromNickname string                  `json:"from_nickname"`
	FromAvatar   string                  `json:"from_avatar"`
	Type         models.NotificationType `json:"type"`
	RecordID     uint                    `json:"record_id"`
	RecordTitle  string                  `json:"record_title"`
	CommentID    *uint                   `json:"comment_id"`
	Content      string                  `json:"content"`
	IsRead       bool                    `json:"is_read"`
	CreatedAt    time.Time               `json:"created_at"`
}

// GetUnreadNotifications 获取未读通知数量
func GetUnreadNotifications(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	var count int64
	models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"unread_count": count,
	})
}

// GetNotifications 获取通知列表
func GetNotifications(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var notifications []models.Notification
	var total int64

	query := models.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
	query.Count(&total)

	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应
	responses := make([]NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		var fromUser models.User
		var record models.Record

		response := NotificationResponse{
			ID:         n.ID,
			UserID:     n.UserID,
			FromUserID: n.FromUserID,
			Type:       n.Type,
			RecordID:   n.RecordID,
			CommentID:  n.CommentID,
			Content:    n.Content,
			IsRead:     n.IsRead,
			CreatedAt:  n.CreatedAt,
		}

		// 获取发送者信息
		if err := models.DB.First(&fromUser, n.FromUserID).Error; err == nil {
			response.FromUsername = fromUser.Username
			response.FromNickname = fromUser.Nickname
			if fromUser.Avatar != "" {
				ctx := context.Background()
				signedURL, err := minioClient.GetPresignedURL(ctx, fromUser.Avatar, 7*24*time.Hour)
				if err == nil {
					response.FromAvatar = signedURL
				}
			}
		}

		// 获取记录标题
		if err := models.DB.First(&record, n.RecordID).Error; err == nil {
			response.RecordTitle = record.Title
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": responses,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// MarkNotificationAsRead 标记通知为已读
func MarkNotificationAsRead(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知 ID"})
		return
	}

	var notification models.Notification
	if err := models.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	// 只能标记自己的通知
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作此通知"})
		return
	}

	notification.IsRead = true
	if err := models.DB.Save(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func MarkAllNotificationsAsRead(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	if err := models.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}
