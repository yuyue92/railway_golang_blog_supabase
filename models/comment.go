package models

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content" binding:"required"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	PostID  uint   `gorm:"not null;index" json:"post_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

func (Comment) TableName() string {
	return "blogcomments"
}
