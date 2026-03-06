package handlers

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
)

// CreateGroup 创建群组
func CreateGroup(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	var req struct {
		Name        string `json:"name" binding:"required,min=2,max=100"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := models.Group{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := models.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建群组失败：" + err.Error()})
		return
	}

	// 自动将创建者添加为群组成员
	userGroup := models.UserGroup{
		UserID:  userID,
		GroupID: group.ID,
		Role:    "owner",
	}
	models.DB.Create(&userGroup)

	// 获取成员数量
	var memberCount int64
	models.DB.Model(&models.UserGroup{}).Where("group_id = ?", group.ID).Count(&memberCount)

	c.JSON(http.StatusCreated, group.ToResponse(memberCount))
}

// GetGroups 获取我加入的群组列表
func GetGroups(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	var groups []models.Group
	// 查询用户加入的群组（包括作为 owner 和 member）
	models.DB.Preload("Owner").
		Joins("JOIN user_groups ON user_groups.group_id = groups.id").
		Where("user_groups.user_id = ?", userID).
		Order("groups.created_at DESC").
		Find(&groups)

	result := make([]gin.H, 0, len(groups))
	for _, g := range groups {
		var memberCount int64
		models.DB.Model(&models.UserGroup{}).Where("group_id = ?", g.ID).Count(&memberCount)

		// 生成头像签名 URL（如果有头像）
		avatar := ""
		if g.Avatar != "" {
			ctx := context.Background()
			signedURL, err := minioClient.GetPresignedURL(ctx, g.Avatar, 7*24*time.Hour)
			if err == nil {
				avatar = signedURL
			}
		}

		result = append(result, gin.H{
			"id":           g.ID,
			"name":         g.Name,
			"description":  g.Description,
			"avatar":       avatar,
			"owner_id":     g.OwnerID,
			"owner_name":   g.Owner.Username,
			"member_count": memberCount,
			"created_at":   g.CreatedAt,
			"updated_at":   g.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetGroup 获取群组详情
func GetGroup(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.Preload("Owner").First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 获取成员数量
	var memberCount int64
	models.DB.Model(&models.UserGroup{}).Where("group_id = ?", gid).Count(&memberCount)

	// 生成头像签名 URL（如果有头像）
	groupResp := group.ToResponse(memberCount)
	if group.Avatar != "" {
		ctx := context.Background()
		signedURL, err := minioClient.GetPresignedURL(ctx, group.Avatar, 7*24*time.Hour)
		if err == nil {
			groupResp.Avatar = signedURL
		}
	}

	c.JSON(http.StatusOK, groupResp)
}

// UpdateGroup 更新群组
func UpdateGroup(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 只有群主可以修改群组
	if group.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改群组"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required,min=2,max=100"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group.Name = req.Name
	group.Description = req.Description

	if err := models.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新群组失败：" + err.Error()})
		return
	}

	var memberCount int64
	models.DB.Model(&models.UserGroup{}).Where("group_id = ?", gid).Count(&memberCount)

	c.JSON(http.StatusOK, group.ToResponse(memberCount))
}

// DeleteGroup 删除群组
func DeleteGroup(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 只有群主可以删除群组
	if group.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除群组"})
		return
	}

	if err := models.DB.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除群组失败：" + err.Error()})
		return
	}

	// 删除群组成员关联
	models.DB.Where("group_id = ?", gid).Delete(&models.UserGroup{})

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AddGroupMember 添加群组成员
func AddGroupMember(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 只有群主或管理员可以添加成员
	var userGroup models.UserGroup
	if err := models.DB.Where("user_id = ? AND group_id = ?", userID, gid).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}
	if userGroup.Role != "owner" && userGroup.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限添加成员"})
		return
	}

	var req struct {
		Usernames []string `json:"usernames" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户并添加
	var users []models.User
	models.DB.Where("username IN ?", req.Usernames).Find(&users)

	addedCount := 0
	for _, user := range users {
		// 检查是否已是成员
		var existing models.UserGroup
		if err := models.DB.Where("user_id = ? AND group_id = ?", user.ID, gid).First(&existing).Error; err == nil {
			continue // 已是成员，跳过
		}

		// 添加成员
		member := models.UserGroup{
			UserID:  user.ID,
			GroupID: gid,
			Role:    "member",
		}
		if err := models.DB.Create(&member).Error; err == nil {
			addedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功添加 " + strconv.Itoa(addedCount) + " 名成员"})
}

// RemoveGroupMember 移除群组成员
func RemoveGroupMember(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 只有群主或管理员可以移除成员
	var userGroup models.UserGroup
	if err := models.DB.Where("user_id = ? AND group_id = ?", userID, gid).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}
	if userGroup.Role != "owner" && userGroup.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限移除成员"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 不能移除群主
	if req.UserID == group.OwnerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能移除群主"})
		return
	}

	if err := models.DB.Where("user_id = ? AND group_id = ?", req.UserID, gid).Delete(&models.UserGroup{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除成员失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除成功"})
}

// GetGroupMembers 获取群组成员列表
func GetGroupMembers(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var members []models.GroupMemberResponse
	models.DB.Table("user_groups").
		Select("user_groups.user_id, user_groups.role, user_groups.created_at as joined_at, users.username, users.nickname, users.avatar").
		Joins("JOIN users ON users.id = user_groups.user_id").
		Where("user_groups.group_id = ?", gid).
		Scan(&members)
	// 对头像进行签名
	for i := range members {
		members[i].Avatar, _ = minioClient.GetPresignedURL(context.Background(), members[i].Avatar, 7*24*time.Hour)
	}
	c.JSON(http.StatusOK, gin.H{"members": members})
}

// LeaveGroup 退出群组
func LeaveGroup(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 群主不能退出，只能转让或解散
	if group.OwnerID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "群主不能退出群组，请先转让群主或解散群组"})
		return
	}

	if err := models.DB.Where("user_id = ? AND group_id = ?", userID, gid).Delete(&models.UserGroup{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出群组失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

// UploadGroupAvatar 上传群组头像
func UploadGroupAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组 ID"})
		return
	}
	gid := uint(groupID)

	var group models.Group
	if err := models.DB.First(&group, gid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 只有群主可以修改群组头像
	if group.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改群组头像"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件"})
		return
	}

	// 检查文件大小（最大 5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 5MB"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValid := false
	for _, e := range validExts {
		if ext == e {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 JPG/PNG/GIF/WEBP 格式"})
		return
	}

	// 生成新文件名
	newFilename := "group_avatar_" + time.Now().Format("20060102150405") + "_" + file.Filename
	objectName := "groups/avatars/" + newFilename

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
	contentType := "image/jpeg"
	if ext == ".png" {
		contentType = "image/png"
	} else if ext == ".gif" {
		contentType = "image/gif"
	} else if ext == ".webp" {
		contentType = "image/webp"
	}

	// 上传到 MinIO
	objectPath, err := minioClient.UploadFile(context.Background(), objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败：" + err.Error()})
		return
	}

	// 只保存相对路径到数据库
	avatarPath := objectPath

	// 更新群组头像（只保存相对路径）
	group.Avatar = avatarPath
	if err := models.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败"})
		return
	}

	// 获取成员数量
	var memberCount int64
	models.DB.Model(&models.UserGroup{}).Where("group_id = ?", gid).Count(&memberCount)

	// 生成签名 URL（7 天有效期）
	ctx := context.Background()
	signedURL, err := minioClient.GetPresignedURL(ctx, avatarPath, 7*24*time.Hour)
	if err != nil {
		log.Printf("生成签名 URL 失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成签名 URL 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"avatar":  signedURL,
		"group":   group.ToResponse(memberCount),
	})
}

// SearchUsers 搜索用户（用于添加成员）
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	var users []models.User
	models.DB.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(20).
		Find(&users)

	result := make([]gin.H, 0, len(users))
	for _, u := range users {
		result = append(result, gin.H{
			"id":       u.ID,
			"username": u.Username,
			"nickname": u.Nickname,
		})
	}

	c.JSON(http.StatusOK, result)
}

// parseUsernames 解析用户名列表
func parseUsernames(usernames string) []string {
	if usernames == "" {
		return []string{}
	}
	list := strings.Split(usernames, ",")
	result := make([]string, 0, len(list))
	for _, u := range list {
		u = strings.TrimSpace(u)
		if u != "" {
			result = append(result, u)
		}
	}
	return result
}
