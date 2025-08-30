// package service ...
package service

import (
	"errors"

	"github.com/KeLes-Coding/gopress/internal/dao"
	"github.com/KeLes-Coding/gopress/internal/model"
	"gorm.io/gorm"
)

// PostService 结构体封装了所有与文章相关的业务逻辑。
type PostService struct{}

// NewPostService 是 PostService 的工厂函数。
func NewPostService() *PostService {
	return &PostService{}
}

// CreatePostDTO (Data Transfer Object) 用于封装创建文章时需要的所有数据。
type CreatePostDTO struct {
	Title      string
	Content    string
	Summary    string
	Status     int
	UserID     uint
	CategoryID uint
	TagIDs     []uint // 标签 ID 列表
}

// Create 用于创建一篇新文章。
func (s *PostService) Create(dto *CreatePostDTO) (*model.Post, error) {
	db := dao.GetDB()
	var category model.Category
	var tags []model.Tag

	// 声明一个 newPost 变量，用于在事务内外传递数据
	newPost := &model.Post{
		Title:      dto.Title,
		Content:    dto.Content,
		Summary:    dto.Summary,
		Status:     dto.Status,
		UserID:     dto.UserID,
		CategoryID: dto.CategoryID,
	}

	// 使用事务 (Transaction) 来确保数据一致性。
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. 校验 CategoryID 是否有效
		if err := tx.First(&category, dto.CategoryID).Error; err != nil {
			return errors.New("无效的分类 ID")
		}

		// 2. 校验所有 TagID 是否有效
		if len(dto.TagIDs) > 0 {
			var count int64
			if err := tx.Model(&model.Tag{}).Where("id IN ?", dto.TagIDs).Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(dto.TagIDs)) {
				return errors.New("包含无效的标签 ID")
			}
			// 获取标签对象以便于关联
			if err := tx.Find(&tags, dto.TagIDs).Error; err != nil {
				return err
			}
			// 将查询到的 tag 实例赋给 newPost
			newPost.Tags = tags
		}

		// 3. 创建 Post
		// 在事务中创建 post 记录
		if err := tx.Create(newPost).Error; err != nil {
			return err
		}

		// 事务成功，返回 nil
		return nil
	})

	if err != nil {
		return nil, err
	}

	// --- 错误修正 ---
	// 在事务成功后, GORM 会自动将新创建记录的 ID 回填到 newPost.ID 字段中。
	// 我们现在可以直接使用 newPost.ID 来查询完整的、预加载了关联数据的文章。
	// 不再需要一个未定义的 lastInsertId 变量。
	var createdPost model.Post
	if err := db.Preload("User").Preload("Category").Preload("Tags").First(&createdPost, newPost.ID).Error; err != nil {
		return nil, err
	}

	return &createdPost, nil
}

// ListPostsDTO 封装了查询文章列表时的参数。
type ListPostsDTO struct {
	Page     int // 页码
	PageSize int // 每页数量
}

// ListResponseDTO 封装了文章列表和总数，用于返回给上层。
type ListResponseDTO struct {
	Posts      []model.Post `json:"post"`
	TotalCount int64        `json:"total_count"`
}

// List 用于获取文章分页列表。
func (s *PostService) List(dto *ListPostsDTO) (*ListResponseDTO, error) {
	db := dao.GetDB()
	var posts []model.Post
	var totalCount int64

	// 计算 offset
	offset := (dto.Page - 1) * dto.PageSize

	// 查询总数
	if err := db.Model(&model.Post{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 查询分页数据，并预加载关联数据
	if err := db.Preload("User").Preload("Category").Preload("Tags").Order("created_at DESC").Limit(dto.PageSize).Offset(offset).Find(&posts).Error; err != nil {
		return nil, err
	}

	return &ListResponseDTO{
		Posts:      posts,
		TotalCount: totalCount,
	}, nil
}

// GetByID 用于根据 ID 获取单篇文章的详细信息。
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	db := dao.GetDB()
	var post model.Post
	if err := db.Preload("User").Preload("Category").Preload("Tags").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	return &post, nil
}

// UpdatePostDTO 封装了更新文章时需要的所有数据。
type UpdatePostDTO struct {
	ID         uint
	Title      string
	Content    string
	Summary    string
	Status     int
	CategoryID uint
	TagIDs     []uint
}

// Update 用于更新一篇文章。
func (s *PostService) Update(dto *UpdatePostDTO) (*model.Post, error) {
	db := dao.GetDB()
	var post model.Post
	var category model.Category
	var tags []model.Tag

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. 查找要更新的文章是否存在
		if err := tx.First(&post, dto.ID).Error; err != nil {
			return errors.New("文章不存在")
		}

		// 2. 校验 CategoryID 是否有效
		if err := tx.First(&category, dto.CategoryID).Error; err != nil {
			return errors.New("无效的分类 ID")
		}

		// 3. 校验所有 TagID 是否有效
		if len(dto.TagIDs) > 0 {
			var count int64
			if err := tx.Model(&model.Tag{}).Where("id IN ?", dto.TagIDs).Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(dto.TagIDs)) {
				return errors.New("包含无效的标签 ID")
			}
			if err := tx.Find(&tags, dto.TagIDs).Error; err != nil {
				return err
			}
		}

		// 4. 更新文章基本信息
		post.Title = dto.Title
		post.Content = dto.Content
		post.Summary = dto.Summary
		post.Status = dto.Status
		post.CategoryID = dto.CategoryID

		if err := tx.Save(&post).Error; err != nil {
			return err
		}

		// 5. 更新文章与标签的多对多关系
		// GORM 的 Association().Replace() 会删除旧的关联，并添加新的关联。
		if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 重新查询以返回完整的、预加载了所有关联数据的文章
	var updatedPost model.Post
	if err := db.Preload("User").Preload("Category").Preload("Tags").First(&updatedPost, dto.ID).Error; err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

// Delete 用于根据 ID 删除一篇文章。
func (s *PostService) Delete(id uint) error {
	db := dao.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var post model.Post
		// 首先需要查找文章已进行关联删除
		if err := tx.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("文章不存在")
			}
			return err
		}

		// GORM 在删除主记录时，会自动删除其在连接表中的关联记录。
		// 我们需要先清除关联，再删除文章本身
		if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
			return err
		}

		// 删除文章
		if err := tx.Delete(&model.Post{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
