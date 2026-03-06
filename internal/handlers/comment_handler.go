package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	recordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录 ID"})
		return
	}
	rid := uint(recordID)

	// 检查记录是否存在
	var record models.Record
	if err := models.DB.First(&record, rid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 检查可见性权限
	if record.Visibility == models.VisibilityPrivate && (record.UserID == nil || *record.UserID != userID) {
		if record.GroupID != nil {
			userGroupIDs := getUserGroupIDs(userID)
			hasPermission := false
			for _, gid := range userGroupIDs {
				if gid == *record.GroupID {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权限评论此记录"})
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限评论此记录"})
			return
		}
	}

	var req struct {
		Content  string `json:"content" binding:"required,min=1,max=500"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 计算楼层号
	floor := 1
	if req.ParentID != nil {
		// 回复评论，获取父评论的楼层
		var parent models.Comment
		if err := models.DB.First(&parent, *req.ParentID).Error; err == nil {
			floor = parent.Floor
		}
	} else {
		// 主评论，获取最大楼层 +1
		var maxFloor models.Comment
		if err := models.DB.Where("record_id = ? AND parent_id IS NULL").Order("floor DESC").First(&maxFloor).Error; err == nil {
			floor = maxFloor.Floor + 1
		}
	}

	comment := models.Comment{
		RecordID: rid,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Floor:    floor,
	}

	if err := models.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败：" + err.Error()})
		return
	}

	// 发送通知
	recordOwnerID := record.UserID
	if recordOwnerID != nil && *recordOwnerID != userID {
		// 评论的是别人的记录，发送通知
		notification := models.Notification{
			UserID:     *recordOwnerID,
			FromUserID: userID,
			Type:       models.NotificationTypeComment,
			RecordID:   rid,
			CommentID:  &comment.ID,
			Content:    req.Content,
		}
		models.DB.Create(&notification)
	}

	// 如果是回复评论，给被回复者发送通知
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := models.DB.First(&parentComment, *req.ParentID).Error; err == nil {
			if parentComment.UserID != userID {
				notification := models.Notification{
					UserID:     parentComment.UserID,
					FromUserID: userID,
					Type:       models.NotificationTypeReply,
					RecordID:   rid,
					CommentID:  &comment.ID,
					Content:    req.Content,
				}
				models.DB.Create(&notification)
			}
		}
	}

	// 查询评论详情
	var result models.Comment
	models.DB.Preload("User").First(&result, comment.ID)

	// 生成头像签名 URL
	commentResp := result.ToResponse(0)
	if result.User != nil && result.User.Avatar != "" {
		ctx := context.Background()
		signedURL, err := minioClient.GetPresignedURL(ctx, result.User.Avatar, 7*24*time.Hour)
		if err == nil {
			commentResp.Avatar = signedURL
		}
	}

	c.JSON(http.StatusCreated, commentResp)
}

// GetComments 获取记录的评论列表（按楼层分组）
func GetComments(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录 ID"})
		return
	}
	rid := uint(recordID)

	// 检查记录是否存在
	var record models.Record
	if err := models.DB.First(&record, recordID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	userID := uint(0)
	if exists {
		userID = userIDInterface.(uint)
	}

	// 检查可见性权限
	if record.Visibility == models.VisibilityPrivate && (record.UserID == nil || *record.UserID != userID) {
		if record.GroupID != nil {
			userGroupIDs := getUserGroupIDs(userID)
			hasPermission := false
			for _, gid := range userGroupIDs {
				if gid == *record.GroupID {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此记录评论"})
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此记录评论"})
			return
		}
	}

	// 获取所有主评论（parent_id IS NULL）
	var comments []models.Comment
	models.DB.Preload("User").
		Where("record_id = ? AND parent_id IS NULL", rid).
		Order("floor ASC").
		Find(&comments)

	result := make([]models.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		// 获取回复数量
		var replyCount int64
		models.DB.Model(&models.Comment{}).
			Where("record_id = ? AND parent_id = ?", rid, comment.ID).
			Count(&replyCount)

		// 获取回复列表（最多 5 条）
		var replies []models.Comment
		models.DB.Preload("User").
			Where("record_id = ? AND parent_id = ?", rid, comment.ID).
			Order("created_at ASC").
			Limit(5).
			Find(&replies)

		commentResp := comment.ToResponse(replyCount)

		// 生成头像签名 URL
		if comment.User != nil && comment.User.Avatar != "" {
			ctx := context.Background()
			signedURL, err := minioClient.GetPresignedURL(ctx, comment.User.Avatar, 7*24*time.Hour)
			if err == nil {
				commentResp.Avatar = signedURL
			}
		}

		// 处理回复
		replyResps := make([]models.CommentResponse, 0, len(replies))
		for _, reply := range replies {
			replyResp := reply.ToResponse(0)
			if reply.User != nil && reply.User.Avatar != "" {
				ctx := context.Background()
				signedURL, err := minioClient.GetPresignedURL(ctx, reply.User.Avatar, 7*24*time.Hour)
				if err == nil {
					replyResp.Avatar = signedURL
				}
			}
			replyResps = append(replyResps, replyResp)
		}
		commentResp.Replies = replyResps

		result = append(result, commentResp)
	}

	c.JSON(http.StatusOK, gin.H{"comments": result, "total": len(result)})
}

// GetCommentReplies 获取评论的回复列表
func GetCommentReplies(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录 ID"})
		return
	}
	rid := uint(recordID)

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论 ID"})
		return
	}
	cid := uint(commentID)

	// 获取回复列表
	var replies []models.Comment
	models.DB.Preload("User").
		Where("record_id = ? AND parent_id = ?", rid, cid).
		Order("created_at ASC").
		Find(&replies)

	result := make([]models.CommentResponse, 0, len(replies))
	for _, reply := range replies {
		replyResp := reply.ToResponse(0)
		if reply.User != nil && reply.User.Avatar != "" {
			ctx := context.Background()
			signedURL, err := minioClient.GetPresignedURL(ctx, reply.User.Avatar, 7*24*time.Hour)
			if err == nil {
				replyResp.Avatar = signedURL
			}
		}
		result = append(result, replyResp)
	}

	c.JSON(http.StatusOK, gin.H{"replies": result})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论 ID"})
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限：只有评论者本人或记录所有者可以删除
	if comment.UserID != userID {
		var record models.Record
		if err := models.DB.First(&record, comment.RecordID).Error; err == nil {
			if record.UserID == nil || *record.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此评论"})
				return
			}
		}
	}

	if err := models.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论 ID"})
		return
	}

	// 点赞数 +1
	if err := models.DB.Model(&models.Comment{}).
		Where("id = ?", commentID).
		Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	// 获取更新后的点赞数
	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"like_count": comment.LikeCount})
}

// GetMyComments 获取我的评论列表
func GetMyComments(c *gin.Context) {
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

	var comments []models.Comment
	var total int64

	models.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&total)
	models.DB.Preload("User").Preload("Record").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&comments)

	result := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		item := gin.H{
			"id":         comment.ID,
			"record_id":  comment.RecordID,
			"content":    comment.Content,
			"floor":      comment.Floor,
			"like_count": comment.LikeCount,
			"created_at": comment.CreatedAt,
		}
		if comment.Record != nil {
			item["record_title"] = comment.Record.Title
			item["record_type"] = comment.Record.Type
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":  result,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
