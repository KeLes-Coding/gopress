package model

import "time"

// Post 模型定义了文章的数据结构。
// 它将映射到数据库中的 `posts` 表。
type Post struct {
	ID      uint   `gorm:"primarykey"`
	Title   string `gorm:"type:varchar(255);not null"` // 文章标题
	Content string `gorm:"type:longtext;not null"`     // 文章内容，使用 longtext 以存储较长的文本
	Summary string `gorm:"type:text"`                  // 文章摘要
	Status  int    `gorm:"type:tinyint;default:1"`     // 状态 (0:草稿, 1:发布)

	// --- 关联字段 ---

	// UserID 是一个外键，关联到 User 模型的 ID。
	// GORM 会根据字段名 `UserID` 自动推断它与 `User` 模型的关系。
	UserID uint `gorm:"not null"`
	// User 字段代表了“属于”(Belongs To)关系。
	// GORM 在查询 Post 时，可以通过 Preload("User") 来自动填充这个字段。
	// `gorm:"foreignKey:UserID"` 明确指定了用于此关系的外键。
	User User `gorm:"foreignKey:UserID"`

	// CategoryID 是一个外键，关联到 Category 模型的 ID。
	CategoryID uint `gorm:"not null"`
	// Category 字段代表了“属于”关系。
	Category Category `gorm:"foreignKey:CategoryID"`

	// Tags 字段代表了“多对多”(Many To Many)关系。
	// `gorm:"many2many:post_tags"` 指定了连接表的名称为 `post_tags`。
	// GORM 会自动创建并管理这个连接表。
	Tags []Tag `gorm:"many2many:post_tags"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 方法用于显式指定模型对应的数据库表名。
func (Post) TableName() string {
	return "posts"
}
