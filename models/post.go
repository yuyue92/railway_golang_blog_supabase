package models

import "gorm.io/gorm"

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"not null;size:200" json:"title" binding:"required"`
	Content  string    `gorm:"type:text;not null" json:"content" binding:"required"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

func (Post) TableName() string {
	return "blogposts"
}
