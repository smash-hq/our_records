package models

import (
	"time"

	"gorm.io/gorm"
)

// Group 群组模型
type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(500)" json:"description"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members     []User         `gorm:"many2many:user_groups;" json:"members,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UserGroup 用户群组关联模型
type UserGroup struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	GroupID   uint      `gorm:"primaryKey;autoIncrement:false" json:"group_id"`
	Role      string    `gorm:"type:varchar(20);default:'member'" json:"role"` // owner, admin, member
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Group     Group     `gorm:"foreignKey:GroupID" json:"-"`
}

// GroupResponse 群组响应
type GroupResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	OwnerID     uint      `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	MemberCount int64     `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse 转换为响应对象
func (g *Group) ToResponse(memberCount int64) GroupResponse {
	ownerName := ""
	if g.Owner != nil {
		ownerName = g.Owner.Username
	}
	return GroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Avatar:      g.Avatar,
		OwnerID:     g.OwnerID,
		OwnerName:   ownerName,
		MemberCount: memberCount,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

// GroupMemberResponse 群组成员响应
type GroupMemberResponse struct {
	UserID   uint      `json:"user_id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
