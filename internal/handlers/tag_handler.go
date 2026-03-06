package handlers

import (
	"net/http"
	"strings"

	"our_records/internal/models"

	"github.com/gin-gonic/gin"
)

// GetMyTags 获取用户所有记录中使用过的标签
func GetMyTags(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	// 获取用户的所有记录
	var records []models.Record
	if err := models.DB.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 提取所有标签
	tagSet := make(map[string]bool)
	for _, record := range records {
		if record.Tags != "" {
			tags := strings.Split(record.Tags, ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tagSet[tag] = true
				}
			}
		}
	}

	// 转换为列表
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}
