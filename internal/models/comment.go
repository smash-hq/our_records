package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RecordID  uint           `gorm:"not null;index" json:"record_id"`
	Record    *Record        `gorm:"foreignKey:RecordID" json:"-"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`            // 父评论 ID，用于回复
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"-"`      // 父评论
	Content   string         `gorm:"type:text;not null" json:"content"` // 评论内容
	Floor     int            `gorm:"not null" json:"floor"`             // 楼层号
	LikeCount int            `gorm:"default:0" json:"like_count"`       // 点赞数
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID         uint              `json:"id"`
	RecordID   uint              `json:"record_id"`
	UserID     uint              `json:"user_id"`
	Username   string            `json:"username"`
	Nickname   string            `json:"nickname"`
	Avatar     string            `json:"avatar"`
	ParentID   *uint             `json:"parent_id"`
	Content    string            `json:"content"`
	Floor      int               `json:"floor"`
	LikeCount  int               `json:"like_count"`
	ReplyCount int64             `json:"reply_count"` // 回复数量
	CreatedAt  time.Time         `json:"created_at"`
	Replies    []CommentResponse `json:"replies,omitempty"` // 子评论（回复）
}

// ToResponse 转换为响应对象
func (c *Comment) ToResponse(replyCount int64) CommentResponse {
	username := ""
	nickname := ""
	avatar := ""
	if c.User != nil {
		username = c.User.Username
		nickname = c.User.Nickname
		avatar = c.User.Avatar
	}
	return CommentResponse{
		ID:         c.ID,
		RecordID:   c.RecordID,
		UserID:     c.UserID,
		Username:   username,
		Nickname:   nickname,
		Avatar:     avatar,
		ParentID:   c.ParentID,
		Content:    c.Content,
		Floor:      c.Floor,
		LikeCount:  c.LikeCount,
		ReplyCount: replyCount,
		CreatedAt:  c.CreatedAt,
	}
}
