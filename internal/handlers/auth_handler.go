package handlers

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"our_records/internal/middleware"
	"our_records/internal/models"
	minioClient "our_records/pkg/minio"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Nickname: req.Nickname,
	}
	if req.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "user": user.ToResponse()})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 JWT Token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 转换为响应对象并生成头像签名 URL
	userResp := user.ToResponse()
	if user.Avatar != "" {
		ctx := context.Background()
		signedURL, err := minioClient.GetPresignedURL(ctx, user.Avatar, 7*24*time.Hour)
		if err == nil {
			userResp.Avatar = signedURL
		}
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  userResp,
	})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 生成头像签名 URL（如果有头像）
	userResp := user.ToResponse()
	if user.Avatar != "" {
		ctx := context.Background()
		signedURL, err := minioClient.GetPresignedURL(ctx, user.Avatar, 7*24*time.Hour)
		if err == nil {
			userResp.Avatar = signedURL
		}
	}

	c.JSON(http.StatusOK, userResp)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
		return
	}

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// UploadAvatar 上传用户头像
func UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDInterface.(uint)

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
	newFilename := "avatar_" + time.Now().Format("20060102150405") + "_" + file.Filename
	objectName := "avatars/" + newFilename

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

	// 更新用户头像（只保存相对路径）
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 删除旧头像（如果有）
	if user.Avatar != "" && !strings.Contains(user.Avatar, "default") {
		// 从路径中提取对象名，删除旧文件
		oldObjectPath := user.Avatar
		ctx := context.Background()
		minioClient.DeleteFile(ctx, oldObjectPath)
	}

	user.Avatar = avatarPath
	if err := models.DB.Save(&user).Error; err != nil {
		log.Printf("更新头像失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败"})
		return
	}

	log.Printf("头像上传成功：user_id=%d, avatar_path=%s", userID, avatarPath)

	// 重新查询用户，生成签名 URL
	var updatedUser models.User
	if err := models.DB.First(&updatedUser, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	// 生成签名 URL（7 天有效期）
	ctx := context.Background()
	signedURL, err := minioClient.GetPresignedURL(ctx, avatarPath, 7*24*time.Hour)
	if err != nil {
		log.Printf("生成签名 URL 失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成签名 URL 失败"})
		return
	}

	// 构建响应
	userResp := updatedUser.ToResponse()
	userResp.Avatar = signedURL

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"avatar":  signedURL,
		"user":    userResp,
	})
}
