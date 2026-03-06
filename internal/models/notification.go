package models

import (
	"time"
)

// NotificationType 通知类型
type NotificationType string

const (
	NotificationTypeComment     NotificationType = "comment"      // 评论我的记录
	NotificationTypeReply       NotificationType = "reply"        // 回复我的评论
	NotificationTypeRecordLike  NotificationType = "record_like"  // 点赞我的记录（预留）
	NotificationTypeCommentLike NotificationType = "comment_like" // 点赞我的评论（预留）
)

// Notification 通知模型
type Notification struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	UserID     uint             `gorm:"not null;index" json:"user_id"` // 接收通知的用户
	FromUserID uint             `gorm:"not null" json:"from_user_id"`  // 发送通知的用户
	Type       NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	RecordID   uint             `gorm:"not null" json:"record_id"`
	CommentID  *uint            `json:"comment_id"`
	Content    string           `gorm:"type:text" json:"content"` // 通知内容（评论/回复的内容）
	IsRead     bool             `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time        `json:"created_at"`
}
