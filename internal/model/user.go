// package model 存放所有与数据库表结构对应的 Go 结构体。
// 这些结构体被称为“模型(Model)”，是 ORM (Object-Relational Mapping) 的核心。
package model

import "time"

// User 模型定义了用户的数据结构，它将映射到数据库中的 `users` 表。
// GORM 会自动将结构体名 `User` 转换为蛇形复数 `users` 作为表名。
// 我们也可以通过实现 TableName() 方法来显式指定表名。
type User struct {
	// GORM 标签(tag)用于定义字段在数据库中的属性。

	// `gorm:"primarykey"`: 将此字段设置为主键。
	ID uint `gorm:"primarykey"`

	// `gorm:"type:varchar(50);unique;not null"`:
	// - type:varchar(50): 指定数据库中的列类型为 VARCHAR，长度为 50。
	// - unique:           为此列添加唯一索引。
	// - not null:         此列不允许为 NULL。
	Username     string `gorm:"type:varchar(50);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(50)"`
	Email        string `gorm:"type:varchar(100);unique"`

	// `gorm:"type:tinyint;default:1"`:
	// - type:tinyint:     指定列类型为 TINYINT。
	// - default:1:        设置此列的默认值为 1。
	Role int `gorm:"type:tinyint;default:1"` // 角色 (0:Admin, 1:User)

	// GORM 的约定：
	// `CreatedAt` 字段: GORM 在创建记录时会自动填充当前时间。
	// `UpdatedAt` 字段: GORM 在创建或更新记录时会自动填充当前时间。
	// 我们也可以使用标签 `gorm:"autoCreateTime"` 和 `gorm:"autoUpdateTime"` 来显式声明。
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 方法允许我们自定义模型对应的数据库表名。
// 如果不实现这个方法，GORM 会默认使用结构体名的蛇形复数形式（例如 `User` -> `users`）。
// 显式地定义表名是一个好习惯，可以避免因命名规则不一致导致的问题。
func (User) TableName() string {
	return "users"
}
