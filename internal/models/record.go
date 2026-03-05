package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// RecordType 记录类型
type RecordType string

const (
	RecordTypeText  RecordType = "text"
	RecordTypeImage RecordType = "image"
	RecordTypeAudio RecordType = "audio"
	RecordTypeVideo RecordType = "video"
)

// Record 记录模型
type Record struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      RecordType     `gorm:"type:varchar(20);not null" json:"type"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	MediaPath string         `gorm:"type:text" json:"media_paths,omitempty"`
	Tags      string         `gorm:"type:varchar(500)" json:"tags"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// GetMediaPaths 获取媒体路径列表
func (r *Record) GetMediaPaths() []string {
	if r.MediaPath == "" {
		return []string{}
	}
	paths := strings.Split(r.MediaPath, ",")
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
