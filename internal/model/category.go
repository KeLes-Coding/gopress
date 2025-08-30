package model

import "time"

// Category 模型定义了文章分类的数据结构。
// 它将映射到数据库中的 `categories` 表。
type Category struct {
	// `gorm:"primarykey"`: 将此字段设置为主键。
	ID uint `gorm:"primarykey"`

	// `gorm:"type:varchar(100);unique;not null"`:
	// - type:varchar(100): 数据库列类型为 VARCHAR，长度 100。
	// - unique:            为此列添加唯一索引，确保分类名不重复。
	// - not null:          此列不允许为 NULL。
	Name string `gorm:"type:varchar(100);unique;not null"`

	// GORM 会在创建记录时自动填充当前时间。
	CreatedAt time.Time
	// GORM 会在创建或更新记录时自动填充当前时间。
	UpdatedAt time.Time
}

// TableName 方法用于显式指定模型对应的数据库表名。
func (Category) TableName() string {
	return "categories"
}
