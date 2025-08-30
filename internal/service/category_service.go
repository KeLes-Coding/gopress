// package service 存放项目的业务逻辑。
package service

import (
	"errors"
	"strings"

	"github.com/KeLes-Coding/gopress/internal/dao"
	"github.com/KeLes-Coding/gopress/internal/model"
	"gorm.io/gorm"
)

// CategoryService 结构体封装了所有与分类相关的业务逻辑。
type CategoryService struct{}

// NewCategoryService 是 CategoryService 的工厂函数。
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// Create 用于创建一个新的分类。
func (s *CategoryService) Create(name string) (*model.Category, error) {
	// 对名称进行基本的处理，例如去除首尾空格
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 检查同名分类是否已存在
	db := dao.GetDB()
	var existingCategory model.Category
	// GORM 的 First 方法在找到记录时返回 nil 错误，未找到时返回 gorm.ErrRecordNotFound
	if err := db.Where("name = ?", trimmedName).First(&existingCategory).Error; err == nil {
		// 如果 err 为 nil，说明已存在同名分类
		return nil, errors.New("分类名称已存在")
	}

	// 创建新的分类实例
	newCategory := &model.Category{Name: trimmedName}

	// 存入数据库
	if err := db.Create(newCategory).Error; err != nil {
		return nil, err
	}

	return newCategory, nil
}

// List 用于获取所有分类的列表。
// 未来可以扩展此方法以支持分页。
func (s *CategoryService) List() ([]model.Category, error) {
	db := dao.GetDB()
	var categories []model.Category
	// GORM 的 Find 方法用于查询多条记录
	// Order("created_at DESC") 表示按创建时间降序排序，最新的分类会排在最前面。
	if err := db.Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 用于更新一个已存在的分类。
// 它需要分类的 ID 和新的名称作为参数。
func (s *CategoryService) Update(id uint, name string) (*model.Category, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, errors.New("分类名称不能为空")
	}

	db := dao.GetDB()

	// 1. 首先，根据 ID 查找分类是否存在
	var category model.Category
	if err := db.First(&category, id).Error; err != nil {
		// 如果 GORM 返回 ErrRecordNotFound，说明该分类不存在。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该分类不存在")
		}
		// 其他数据库错误
		return nil, err
	}

	// 2. 检查新的名称是否存在其他分类的名称冲突
	var existingCategory model.Category
	if err := db.Where("name = ? AND id != ?", trimmedName, id).First(&existingCategory).Error; err == nil {
		return nil, errors.New("该分类名称已存在")
	}

	// 3. 更新分类名称
	category.Name = trimmedName
	if err := db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// Delete 用于根据 ID 删除一个分类。
func (s *CategoryService) Delete(id uint) error {
	db := dao.GetDB()

	// GORM 的 Delete 方法可以通过主键删除记录
	// 它会返回一个 result 对象，我们可以通过 .RowsAffected 检查是否有记录被删除。
	result := db.Delete(&model.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	// 如果没有行受到影响，说明该 ID 的分类原本就不存在
	if result.RowsAffected == 0 {
		return errors.New("该分类不存在")
	}

	return nil
}
