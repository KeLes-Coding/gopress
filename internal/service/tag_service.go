// package service 存放项目的业务逻辑。
package service

import (
	"errors"
	"strings"

	"github.com/KeLes-Coding/gopress/internal/dao"
	"github.com/KeLes-Coding/gopress/internal/model"
	"gorm.io/gorm"
)

// TagService 结构体封装了所有与标签相关的业务逻辑。
type TagService struct{}

// NewTagService 是 TagService 的工厂函数。
func NewTagService() *TagService {
	return &TagService{}
}

// Create 用于创建一个新的标签。
func (s *TagService) Create(name string) (*model.Tag, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, errors.New("抱歉名称不能为空")
	}

	db := dao.GetDB()
	var existingTag model.Tag
	if err := db.Where("name = ?", trimmedName).First(&existingTag).Error; err == nil {
		return nil, errors.New("该标签名称已存在")
	}

	newTag := &model.Tag{Name: trimmedName}
	if err := db.Create(newTag).Error; err != nil {
		return nil, err
	}
	return newTag, nil
}

// List 用于获取所有标签的列表。
func (s *TagService) List() ([]model.Tag, error) {
	db := dao.GetDB()
	var tags []model.Tag
	if err := db.Order("created_at DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Update 用于更新一个已存在的标签。
func (s *TagService) Update(id uint, name string) (*model.Tag, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, errors.New("标签名称不能为空")
	}

	db := dao.GetDB()
	var tag model.Tag
	if err := db.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该标签不存在")
		}
		return nil, err
	}

	var existingTag model.Tag
	if err := db.Where("name = ? AND id != ?", trimmedName, id).First(&existingTag).Error; err == nil {
		return nil, errors.New("该标签名称已存在")
	}

	tag.Name = trimmedName
	if err := db.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Delete 用于根据 ID 删除一个标签。
func (s *TagService) Delete(id uint) error {
	db := dao.GetDB()
	result := db.Delete(&model.Tag{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("该标签不存在")
	}
	return nil
}
